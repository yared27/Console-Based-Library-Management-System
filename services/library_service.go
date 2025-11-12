package services

import (
	"errors"
	"fmt"
	"library_management/models"
	"sync"
	"library_management/concurrency"
)

// type ReservationRequest struct {
// 	BookID   int
// 	MemberID int
// }

type LibraryManager interface {
	AddBook(book models.Book)
	RemoveBook(bookID int)
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) []models.Book
	ReserveBook(bookID, memberID int) error
	AddMember(memberID int,Name string) 
}

type Library struct {
	Books   map[int]models.Book
	Members map[int]models.Member

	mu      sync.Mutex
	Reserve chan concurrency.ReservationRequest
}

func NewLibrary() *Library {
	l := &Library{
		Books:   make(map[int]models.Book),
		Members: make(map[int]models.Member),
		Reserve: make(chan concurrency.ReservationRequest, 100),
	}
	concurrency.StartReservationWorker(l.Reserve, l.Books, &l.mu)
	return l
}

// addbook adds a new book to the library
func (l *Library) AddBook(book models.Book) {
	l.Books[book.ID] = book
}

// remove book from the library
func (l *Library) RemoveBook(bookID int) {
	_, exists := l.Books[bookID]
	if exists {
		delete(l.Books, bookID)
		fmt.Println("Book removed successfully")
	} else {
		fmt.Println("Book not found!")
	}
}

func (l *Library) BorrowBook(bookID int, memberID int) error {
	book, bookExists := l.Books[bookID]
	if !bookExists {
		return errors.New("book not found")
	}

	if book.Status == "Borrowed" {
		return errors.New("book is already borrowed")
	}
	member, memberExists := l.Members[memberID]
	if !memberExists {
		return errors.New("member not found")
	}

	book.Status = "Borrowed"
	l.Books[bookID] = book
	member.BorrowedBooks = append(member.BorrowedBooks, book)
	l.Members[memberID] = member
	return nil
}

func (l *Library) ReturnBook(bookID int, memberID int) error {
	book, bookExists := l.Books[bookID]
	if !bookExists {
		return errors.New("book not found")
	}
	member, memberExists := l.Members[memberID]
	if !memberExists {
		return errors.New("member not found")
	}

	// Remove book from member's borrowed list
	found := false
	for i, b := range member.BorrowedBooks {
		if b.ID == bookID {
			member.BorrowedBooks = append(member.BorrowedBooks[:i], member.BorrowedBooks[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return errors.New("book not borrowed by member")
	}

	// Mark the book as available again and update records
	book.Status = "Available"
	l.Books[bookID] = book
	l.Members[memberID] = member
	return nil
}

func (l *Library) ListAvailableBooks() []models.Book {
	available := []models.Book{}

	for _, book := range l.Books {
		if book.Status == "Available" {
			available = append(available, book)
		}
	}
	return available

}

func (l *Library) ListBorrowedBooks(memberID int) []models.Book {
	member, exists := l.Members[memberID]
	if !exists {
		fmt.Println("Member not found!")
		return nil
	}
	return member.BorrowedBooks
}


func (l *Library) ReserveBook(bookID, memberID int) error {
	if _, exists := l.Books[bookID]; !exists {
		return fmt.Errorf("book with ID %d does not exist", bookID)
	}

	if _, exists := l.Members[memberID]; !exists {
		return fmt.Errorf("Memeber with ID %d does not exists", memberID)
	}

	// send reservation request to the channel (non-blocking)

	go func() {
		l.Reserve <- concurrency.ReservationRequest{BookID: bookID, MemberID: memberID}
	}()
	return nil
}

func (l*Library) AddMember(memberID int,name string){
	l.Members[memberID] = models.Member{
		ID: memberID,
		Name: name,
	}
	fmt.Printf("Member '%s' added with ID %d\n", name, memberID)
}