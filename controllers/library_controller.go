package controllers

import (
	"bufio"
	"fmt"
	"library_management/models"
	"library_management/services"
	"os"
	"strings"
)
func ReciveIDs()(int,int){
			var bookid, memberid int
			fmt.Print("Book ID to borrow: ")
			fmt.Scanln(&bookid)
			fmt.Print("Member ID: ")
			fmt.Scanln(&memberid)
			return bookid,memberid
}

func StartLibraryConsole() {
	library := services.NewLibrary()
	reader := bufio.NewReader(os.Stdin) 
	library.Members[1] = models.Member{ID:1, Name:"Yared"}
	for {
		fmt.Println("\n=== Library Menu ===")
		fmt.Println("1. Add Book")
		fmt.Println("2. Remove Book")
		fmt.Println("3. Borrow Book")
		fmt.Println("4. Return Book")
		fmt.Println("5. List Available Books")
		fmt.Println("6. List Borrowed Books by Member")
		fmt.Println("7. Add New Memeber")
		fmt.Println("8. Reserve Book")
		fmt.Println("9. Exit")
		fmt.Println("Choose an option: ")

		var choice int
		fmt.Scanln(&choice)

		switch choice{
		case 1:
			var id int 
			fmt.Print("Book ID: ")
			fmt.Scanln(&id)
			fmt.Print("Title: ")
			title,_ := reader.ReadString('\n')
			title = strings.TrimSpace(title)
			fmt.Print("Author: ")
			author, _ := reader.ReadString('\n')
			author = strings.TrimSpace(author)
			book := models.Book{ID: id, Title: title, Author: author, Status: "Available"}
			library.AddBook(book)
			fmt.Println("Book Added successfully!")
		case 2:
			var id int
			fmt.Print("Book ID to remove: ")
			fmt.Scanln(&id)
			library.RemoveBook(id)
			fmt.Println("Book removed successfully!")
		
		case 3:
			bookID, memberID := ReciveIDs()
			err := library.BorrowBook(bookID, memberID)
			if err != nil {
				fmt.Println("Error:", err)
			}else{
				fmt.Println("Book borrowed successfully!")
			}
		case 4:
			 bookID, memberID := ReciveIDs()
			 err := library.ReturnBook(bookID, memberID)
			 if err != nil{
				fmt.Println("Error:", err)
			 }else{
				fmt.Println("Book returned successfully!")
			 }
		case 5:
				fmt.Println("\nAvailable Books:")
				for _, b := range library.ListAvailableBooks(){
					fmt.Printf("ID: %d , Title: %s , Author: %s , Status: %s\n", b.ID, b.Title, b.Author, b.Status)
				}
		case 6:
				var memberID int 
				fmt.Print("Member ID: ")
				fmt.Scanln(&memberID)
				books := library.ListBorrowedBooks(memberID)

				if books != nil{
					fmt.Println("Borrowed Books:")
					for _, b := range books {
						fmt.Printf("ID: %d, Title: %s, Author: %s, Status: %s\n", b.ID,b.Title,b.Author, b.Status)
					}
				}
		case 7:
				var memeberID int
				fmt.Println("Enter MemberID:")
				fmt.Scanln(&memeberID)
				fmt.Println("Enter Member Name")
				name, _:= reader.ReadString('\n')
				name = strings.TrimSpace(name)
				fmt.Scanln(&name)
				library.AddMember(memeberID,name)
				fmt.Println("Member added successfully!")
			
		case 8:
				fmt.Println("Reserve Book")
				bookID, memberID := ReciveIDs()
				err := library.ReserveBook(bookID,memberID)
				if err != nil{
					fmt.Println("Error:", err)
				}

		case 9:
				fmt.Println("Exiting... Goodbye!")
				return

		default:
				fmt.Println("Invalid choice, try again")
		}
	}
}