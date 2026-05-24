package scraper

import (
	"context"
	"job-matching-scraper/internal/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// fungsi ini untuk menyimpan hasil scraping kedalam database
func (r *Repository) SaveJobs(ctx context.Context, jobs []entity.Job) (int, int, error) {
	var inserted int
	var duplicated int

	//jumlah data yang akan dimasukkan, agar tidak sekaligus memasukkan banyak data
	batchSize := 50

	for i := 0; i < len(jobs); i += batchSize {
		//menghitung awal index dan akhir index yang akan dipotong dan dikirm ke db
		end := i + batchSize
		//jika end lebih dari jobx, fallback end menjadi jumlah panjang jobs
		if end > len(jobs) {
			end = len(jobs)
		}
		//mengambil data jobs perbatch yang sudah ditentukan
		batch := jobs[i:end]

		//memasukkan jobs kedalam database
		result := r.db.WithContext(ctx).
			//megecek apakah column url terjadi duplikat, jika terjadi jangan lakukan apapun alias skip
			Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "url"}},
				DoNothing: true,
			}).
			Create(&batch)

		if result.Error != nil {
			return inserted, duplicated, result.Error
		}

		//menghitung jumlah insert
		inserted += int(result.RowsAffected)
		//menghitung jumlah dupliat dengan mengurangi panjang bathc dan jumlah rowaffected
		duplicated += len(batch) - int(result.RowsAffected)
	}

	return inserted, duplicated, nil
}
