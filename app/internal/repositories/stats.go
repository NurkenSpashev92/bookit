package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"

	"github.com/nurkenspashev92/bookit/internal/schemas"
)

type StatsRepository struct {
	db *pgxpool.Pool
}

func NewStatsRepository(db *pgxpool.Pool) *StatsRepository {
	return &StatsRepository{db: db}
}

func (r *StatsRepository) GetDashboard(ctx context.Context, ownerID int) (schemas.DashboardStats, error) {
	var s schemas.DashboardStats

	err := r.db.QueryRow(ctx, `
		SELECT
			COUNT(*),
			COALESCE(SUM(view_count), 0),
			COALESCE(SUM(like_count), 0),
			COALESCE(AVG(price), 0),
			COALESCE(AVG(view_count), 0),
			COALESCE(AVG(like_count), 0)
		FROM houses
		WHERE owner_id = $1
	`, ownerID).Scan(&s.TotalHouses, &s.TotalViews, &s.TotalLikes, &s.AveragePrice, &s.AverageViews, &s.AverageLikes)
	if err != nil {
		return s, err
	}

	// Median price
	_ = r.db.QueryRow(ctx, `
		SELECT COALESCE(percentile_cont(0.5) WITHIN GROUP (ORDER BY price), 0)
		FROM houses WHERE owner_id = $1
	`, ownerID).Scan(&s.MedianPrice)

	return s, nil
}

func (r *StatsRepository) GetHouseStats(ctx context.Context, ownerID int) ([]schemas.HouseStatsItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name_en, name_kz, name_ru, slug, price, view_count, like_count
		FROM houses
		WHERE owner_id = $1
		ORDER BY view_count DESC
	`, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []schemas.HouseStatsItem
	for rows.Next() {
		var h schemas.HouseStatsItem
		if err := rows.Scan(&h.ID, &h.NameEN, &h.NameKZ, &h.NameRU, &h.Slug, &h.Price, &h.ViewCount, &h.LikeCount); err != nil {
			return nil, err
		}
		items = append(items, h)
	}
	return items, nil
}

func (r *StatsRepository) GetTopByViews(ctx context.Context, ownerID, limit int) ([]schemas.TopHouseStats, error) {
	return r.getTop(ctx, ownerID, "view_count", limit)
}

func (r *StatsRepository) GetTopByLikes(ctx context.Context, ownerID, limit int) ([]schemas.TopHouseStats, error) {
	return r.getTop(ctx, ownerID, "like_count", limit)
}

func (r *StatsRepository) getTop(ctx context.Context, ownerID int, column string, limit int) ([]schemas.TopHouseStats, error) {
	query := fmt.Sprintf(`
		SELECT id, name_en, name_kz, name_ru, slug, %s
		FROM houses
		WHERE owner_id = $1
		ORDER BY %s DESC
		LIMIT $2
	`, column, column)

	rows, err := r.db.Query(ctx, query, ownerID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []schemas.TopHouseStats
	for rows.Next() {
		var h schemas.TopHouseStats
		if err := rows.Scan(&h.ID, &h.NameEN, &h.NameKZ, &h.NameRU, &h.Slug, &h.Value); err != nil {
			return nil, err
		}
		items = append(items, h)
	}
	return items, nil
}

func (r *StatsRepository) GetLikesOverTime(ctx context.Context, ownerID, days int) ([]schemas.LikesOverTime, error) {
	rows, err := r.db.Query(ctx, `
		SELECT d.date::text, COALESCE(cnt, 0)
		FROM generate_series(
			CURRENT_DATE - make_interval(days => $2),
			CURRENT_DATE,
			'1 day'
		) AS d(date)
		LEFT JOIN (
			SELECT hl.created_at::date AS dt, COUNT(*) AS cnt
			FROM house_likes hl
			INNER JOIN houses h ON h.id = hl.house_id
			WHERE h.owner_id = $1
			  AND hl.created_at >= CURRENT_DATE - make_interval(days => $2)
			GROUP BY dt
		) sub ON sub.dt = d.date
		ORDER BY d.date
	`, ownerID, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []schemas.LikesOverTime
	for rows.Next() {
		var l schemas.LikesOverTime
		if err := rows.Scan(&l.Date, &l.Count); err != nil {
			return nil, err
		}
		items = append(items, l)
	}
	return items, nil
}

func (r *StatsRepository) GetPriceDistribution(ctx context.Context, ownerID int) ([]schemas.PriceDistribution, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			CASE
				WHEN price < 50000 THEN '0-50000'
				WHEN price < 100000 THEN '50000-100000'
				WHEN price < 200000 THEN '100000-200000'
				WHEN price < 500000 THEN '200000-500000'
				ELSE '500000+'
			END AS range,
			COUNT(*) AS count
		FROM houses
		WHERE owner_id = $1
		GROUP BY range
		ORDER BY MIN(price)
	`, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []schemas.PriceDistribution
	for rows.Next() {
		var p schemas.PriceDistribution
		if err := rows.Scan(&p.Range, &p.Count); err != nil {
			return nil, err
		}
		items = append(items, p)
	}
	return items, nil
}

func (r *StatsRepository) GetHouseDetailStats(ctx context.Context, ownerID int, slug string) (schemas.HouseDetailStats, error) {
	var s schemas.HouseDetailStats

	err := r.db.QueryRow(ctx, `
		SELECT id, name_en, name_kz, name_ru, slug, view_count, like_count
		FROM houses WHERE slug=$1 AND owner_id=$2
	`, slug, ownerID).Scan(&s.HouseID, &s.NameEN, &s.NameKZ, &s.NameRU, &s.Slug, &s.TotalViews, &s.TotalLikes)
	if err != nil {
		return s, fmt.Errorf("house not found or not owned by you")
	}

	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return r.db.QueryRow(gctx, `
			SELECT COUNT(DISTINCT COALESCE(user_id::text, ip_address))
			FROM house_views WHERE house_id=$1
		`, s.HouseID).Scan(&s.UniqueViews)
	})

	g.Go(func() error {
		return r.db.QueryRow(gctx, `
			SELECT COUNT(*) FROM bookings WHERE house_id=$1
		`, s.HouseID).Scan(&s.TotalBookings)
	})

	g.Go(func() error {
		rows, err := r.db.Query(gctx, `
			SELECT hv.user_id, COALESCE(u.email, ''), COALESCE(u.first_name, ''), COALESCE(u.last_name, ''),
				   COALESCE(hv.ip_address, ''), hv.created_at
			FROM house_views hv
			LEFT JOIN users u ON u.id = hv.user_id
			WHERE hv.house_id=$1
			ORDER BY hv.created_at DESC LIMIT 50
		`, s.HouseID)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var v schemas.ViewerInfo
			var viewedAt time.Time
			if err := rows.Scan(&v.UserID, &v.Email, &v.FirstName, &v.LastName, &v.IP, &viewedAt); err != nil {
				return err
			}
			v.ViewedAt = viewedAt.Format(time.RFC3339)
			s.RecentViewers = append(s.RecentViewers, v)
		}
		return nil
	})

	g.Go(func() error {
		rows, err := r.db.Query(gctx, `
			SELECT u.id, u.email, COALESCE(u.first_name, ''), COALESCE(u.last_name, ''), hl.created_at
			FROM house_likes hl
			INNER JOIN users u ON u.id = hl.user_id
			WHERE hl.house_id=$1
			ORDER BY hl.created_at DESC LIMIT 50
		`, s.HouseID)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var l schemas.LikerInfo
			var likedAt time.Time
			if err := rows.Scan(&l.UserID, &l.Email, &l.FirstName, &l.LastName, &likedAt); err != nil {
				return err
			}
			l.LikedAt = likedAt.Format(time.RFC3339)
			s.RecentLikers = append(s.RecentLikers, l)
		}
		return nil
	})

	g.Go(func() error {
		rows, err := r.db.Query(gctx, `
			SELECT b.id, u.id, u.email, COALESCE(u.first_name, ''), COALESCE(u.last_name, ''),
				   b.start_date, b.end_date, b.guest_count, b.status, b.total_price, b.created_at
			FROM bookings b
			INNER JOIN users u ON u.id = b.user_id
			WHERE b.house_id=$1
			ORDER BY b.created_at DESC LIMIT 50
		`, s.HouseID)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var b schemas.BookingInfo
			var startDate, endDate, createdAt time.Time
			if err := rows.Scan(&b.ID, &b.UserID, &b.Email, &b.FirstName, &b.LastName,
				&startDate, &endDate, &b.GuestCount, &b.Status, &b.TotalPrice, &createdAt); err != nil {
				return err
			}
			b.StartDate = startDate.Format("2006-01-02")
			b.EndDate = endDate.Format("2006-01-02")
			b.CreatedAt = createdAt.Format(time.RFC3339)
			s.RecentBookings = append(s.RecentBookings, b)
		}
		return nil
	})

	g.Go(func() error {
		rows, err := r.db.Query(gctx, `
			SELECT d.date::text, COALESCE(cnt, 0)
			FROM generate_series(CURRENT_DATE - 29, CURRENT_DATE, '1 day') AS d(date)
			LEFT JOIN (
				SELECT created_at::date AS dt, COUNT(*) AS cnt
				FROM house_views WHERE house_id=$1 AND created_at >= CURRENT_DATE - 29
				GROUP BY dt
			) sub ON sub.dt = d.date ORDER BY d.date
		`, s.HouseID)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var item schemas.LikesOverTime
			if err := rows.Scan(&item.Date, &item.Count); err != nil {
				return err
			}
			s.ViewsPerDay = append(s.ViewsPerDay, item)
		}
		return nil
	})

	g.Go(func() error {
		rows, err := r.db.Query(gctx, `
			SELECT d.date::text, COALESCE(cnt, 0)
			FROM generate_series(CURRENT_DATE - 29, CURRENT_DATE, '1 day') AS d(date)
			LEFT JOIN (
				SELECT created_at::date AS dt, COUNT(*) AS cnt
				FROM house_likes WHERE house_id=$1 AND created_at >= CURRENT_DATE - 29
				GROUP BY dt
			) sub ON sub.dt = d.date ORDER BY d.date
		`, s.HouseID)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var item schemas.LikesOverTime
			if err := rows.Scan(&item.Date, &item.Count); err != nil {
				return err
			}
			s.LikesPerDay = append(s.LikesPerDay, item)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return s, err
	}

	if s.RecentViewers == nil {
		s.RecentViewers = []schemas.ViewerInfo{}
	}
	if s.RecentLikers == nil {
		s.RecentLikers = []schemas.LikerInfo{}
	}
	if s.RecentBookings == nil {
		s.RecentBookings = []schemas.BookingInfo{}
	}
	if s.ViewsPerDay == nil {
		s.ViewsPerDay = []schemas.LikesOverTime{}
	}
	if s.LikesPerDay == nil {
		s.LikesPerDay = []schemas.LikesOverTime{}
	}

	return s, nil
}
