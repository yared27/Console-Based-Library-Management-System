Library Management System – Concurrency Documentation
Overview

This system is an extension of the Library Management System to support concurrent book reservations. It allows multiple members to reserve books at the same time safely, handles automatic reservation expiration, and ensures consistent book availability.

Concurrency Approach
1. Goroutines

Purpose: Execute reservation tasks concurrently without blocking the main menu or other operations.

Implementation:

ReserveBook starts a goroutine to handle reservation expiration after 5 seconds.

Multiple members can attempt reservations simultaneously using separate goroutines.

go func(bookID int) {
    time.Sleep(5 * time.Second)
    // expiration logic here
}(bookID)

2. Mutex (sync.Mutex)

Purpose: Prevent race conditions when updating shared resources (book availability).

Implementation:

Before reading/updating a book’s status, the goroutine locks the mutex.

After updating, it unlocks the mutex using defer.

l.mu.Lock()
defer l.mu.Unlock()
book := l.Books[bookID]
// update book status

3. Channels

Purpose: Queue reservation requests and coordinate asynchronous processing.

Implementation:

A Reserve channel in the Library struct is used to send reservation requests to a worker goroutine (StartReservationWorker).

The worker reads requests from the channel and processes them concurrently.

for req := range l.Reserve {
    // handle reservation
}

4. Auto-Cancellation

If a book is reserved but not borrowed within 5 seconds, a goroutine automatically sets the status back to "Available".

This ensures books are not locked indefinitely and remain available for other members.

Handling Concurrent Requests

Multiple members attempting to reserve the same book will be safely handled:

The first successful reservation locks the book.

Others attempting to reserve the same book receive an error message indicating it is unavailable.

Summary

Goroutines allow asynchronous processing of reservations.

Mutexes ensure thread-safe updates of book data.

Channels manage queued reservation requests efficiently.

The system demonstrates safe concurrent programming in Go while maintaining a responsive console interface.