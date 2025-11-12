package concurrency
import ("fmt"
		"time"
		"library_management/models"
		"sync"
)


type ReservationRequest struct {
	BookID   int
	MemberID int
}


func  StartReservationWorker(reserveChan <- chan ReservationRequest,
	books map[int]models.Book,
	mu *sync.Mutex) {
	go func() {
		for req := range reserveChan {
			mu.Lock()
			book := books[req.BookID]

			if book.Status == "Available" {
				book.Status = "Reserved"
				books[req.BookID] = book
				fmt.Printf("\n[INFO] Book %s reserved by member %d\n", book.Title, req.MemberID)

				go func(bookID int) {
					time.Sleep(5 * time.Second)
					mu.Lock()
					defer mu.Unlock()

					book := books[bookID]

					if book.Status == "Reserved" {
						book.Status = "Available"
						books[bookID] = book
						fmt.Printf("\n[INFO] Reservation for book '%s' expired\n", book.Title)
					}
				}(req.BookID)
			} else {
				fmt.Printf("Book '%s' is not available for reservation\n", book.Title)
			}
			mu.Unlock()
		}

	}()
}
