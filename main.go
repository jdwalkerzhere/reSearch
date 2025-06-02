package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"

	"github.com/jessewalker/reSearch/internal/database"
)

func main() {
	fmt.Println("Welcome to reSearch - AI Research Talent Exploration Tool")

	// Open database connection
	fmt.Println("[DEBUG] Opening database connection...")
	db, err := sql.Open("sqlite3", "research.db")
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
		return
	}
	defer db.Close()
	fmt.Println("[DEBUG] Database connection opened successfully")

	// Test database connection
	err = db.Ping()
	if err != nil {
		fmt.Printf("Error pinging database: %v\n", err)
		return
	}
	fmt.Println("[DEBUG] Database ping successful")

	// Initialize database queries
	fmt.Println("[DEBUG] Initializing database queries...")
	queries := database.New(db)
	fmt.Println("[DEBUG] Database queries initialized")

	// Create a scanner to read user input
	scanner := bufio.NewScanner(os.Stdin)

	// Create context for database operations
	ctx := context.Background()

	// Main application loop
	for {
		fmt.Println("[DEBUG] Displaying main menu")
		displayMainMenu()

		// Get user input
		fmt.Print("\nEnter your choice: ")
		scanner.Scan()
		choice := strings.TrimSpace(scanner.Text())
		fmt.Printf("[DEBUG] User selected option: '%s'\n", choice)

		// Process user selection
		switch choice {
		case "1":
			fmt.Println("\n--- Create New Search ---")
			createNewSearch(ctx, queries, scanner)
			pressEnterToContinue(scanner)
		case "2":
			fmt.Println("\n--- Manage Searches ---")
			manageSearches(ctx, queries, scanner)
			pressEnterToContinue(scanner)
		case "3":
			fmt.Println("\n--- Check New Results ---")
			fmt.Println("This feature will be implemented soon.")
			pressEnterToContinue(scanner)
		case "4":
			fmt.Println("\n--- Fetch Older Results ---")
			fmt.Println("This feature will be implemented soon.")
			pressEnterToContinue(scanner)
		case "5":
			fmt.Println("\n--- Delete Search ---")
			fmt.Println("This feature will be implemented soon.")
			pressEnterToContinue(scanner)
		case "q", "Q", "exit", "quit":
			fmt.Println("Exiting reSearch. Goodbye!")
			return
		default:
			fmt.Println("Invalid option. Please try again.")
			pressEnterToContinue(scanner)
		}
	}
}

// displayMainMenu shows the main menu options to the user
func displayMainMenu() {
	fmt.Println("\n==== Main Menu ====")
	fmt.Println("1. Create New Search")
	fmt.Println("2. Manage Searches")
	fmt.Println("3. Check New Results")
	fmt.Println("4. Fetch Older Results")
	fmt.Println("5. Delete Search")
	fmt.Println("q. Quit")
}

// pressEnterToContinue waits for the user to press Enter
func pressEnterToContinue(scanner *bufio.Scanner) {
	fmt.Print("\nPress Enter to continue...")
	scanner.Scan()
}

// manageSearches allows viewing and managing existing searches
func manageSearches(ctx context.Context, queries *database.Queries, scanner *bufio.Scanner) {
	fmt.Println("[DEBUG] Starting manageSearches function")
	
	// List active searches with pagination
	limit := int64(10)  // Show 10 searches per page
	offset := int64(0)  // Start with the first page
	
	for {
		fmt.Println("\nFetching active searches...")
		listParams := database.ListActiveSearchesParams{
			Limit:  limit,
			Offset: offset,
		}
		
		searches, err := queries.ListActiveSearches(ctx, listParams)
		if err != nil {
			fmt.Printf("Error listing searches: %v\n", err)
			return
		}
		
		if len(searches) == 0 {
			if offset == 0 {
				fmt.Println("No searches found. Create a new search first.")
				return
			} else {
				fmt.Println("No more searches found.")
				// Go back to the previous page
				offset = offset - limit
				if offset < 0 {
					offset = 0
				}
				continue
			}
		}
		
		// Display searches in a formatted table
		fmt.Println("\n+--------+------------------------+----------------------------+---------------+-------------+")
		fmt.Println("| NUMBER | DESCRIPTION            | ARXIV URL                  | LAST FETCH    | ARTICLE CT  |")
		fmt.Println("+--------+------------------------+----------------------------+---------------+-------------+")
		
		for i, search := range searches {
			// Format description (truncate if too long)
			description := search.Description
			if len(description) > 22 {
				description = description[:19] + "..."
			}
			
			// Format arXiv URL (truncate if too long)
			arxivURL := search.ArvixUrl
			if len(arxivURL) > 26 {
				arxivURL = arxivURL[:23] + "..."
			}
			
			// Format last fetch date
			lastFetch := "Never"
			if search.LastFetchDate.Valid {
				lastFetch = search.LastFetchDate.Time.Format("2006-01-02")
			}
			
			fmt.Printf("| %-6d | %-22s | %-26s | %-13s | %-11d |\n", 
				i+1, description, arxivURL, lastFetch, search.ArticleCount)
		}
		fmt.Println("+--------+------------------------+----------------------------+---------------+-------------+")
		
		// Navigation options
		fmt.Println("\nOptions:")
		fmt.Println("  [number] - Select search to view details")
		if offset > 0 {
			fmt.Println("  p - Previous page")
		}
		if len(searches) == int(limit) {
			fmt.Println("  n - Next page")
		}
		fmt.Println("  b - Back to main menu")
		
		fmt.Print("\nEnter choice: ")
		scanner.Scan()
		choice := strings.TrimSpace(scanner.Text())
		fmt.Printf("[DEBUG] Search list - user selected: '%s'\n", choice)
		
		// Handle navigation
		if choice == "b" {
			return
		} else if choice == "n" && len(searches) == int(limit) {
			offset += limit
			continue
		} else if choice == "p" && offset > 0 {
			offset -= limit
			continue
		} else {
			// Try to parse as a number for search selection
			var selectedIndex int
			_, err := fmt.Sscanf(choice, "%d", &selectedIndex)
			if err == nil && selectedIndex > 0 && selectedIndex <= len(searches) {
				viewSearchDetails(ctx, queries, searches[selectedIndex-1].ID, scanner)
			} else {
				// Check if the user entered a letter that's for the details view
				if choice == "d" || choice == "e" || choice == "f" {
					fmt.Println("You need to select a search first (enter its number) before using that option.")
				} else {
					fmt.Println("Invalid selection. Please enter a search number, 'n' (next), 'p' (previous), or 'b' (back).")
				}
				time.Sleep(1 * time.Second) // Brief pause to let user see the message
			}
		}
	}
}

// viewSearchDetails displays detailed information about a specific search
func viewSearchDetails(ctx context.Context, queries *database.Queries, searchID interface{}, scanner *bufio.Scanner) {
	fmt.Printf("[DEBUG] Viewing search details for ID: %v\n", searchID)
	
	// Get detailed search information
	searchStats, err := queries.GetSearchWithStats(ctx, searchID)
	if err != nil {
		fmt.Printf("Error retrieving search details: %v\n", err)
		return
	}
	
	// Display search details
	fmt.Println("\n=== Search Details ===")
	fmt.Printf("ID:              %v\n", searchStats.ID)
	fmt.Printf("Description:     %s\n", searchStats.Description)
	fmt.Printf("arXiv URL:       %s\n", searchStats.ArvixUrl)
	fmt.Printf("Created:         %s\n", searchStats.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Last Updated:    %s\n", searchStats.UpdatedAt.Format("2006-01-02 15:04:05"))
	
	// Format fetch information
	fmt.Printf("Results/Fetch:   %d\n", defaultIfNullInt64(searchStats.ResultsPerFetch, 50))
	
	lastFetch := "Never"
	if searchStats.LastFetchDate.Valid {
		lastFetch = searchStats.LastFetchDate.Time.Format("2006-01-02 15:04:05")
	}
	fmt.Printf("Last Fetch:      %s\n", lastFetch)
	
	// Show statistics
	fmt.Printf("Article Count:   %d\n", searchStats.ArticleCount)
	fmt.Printf("Candidate Count: %d\n", searchStats.CandidateCount)
	
	// Display options
	fmt.Println("\nOptions:")
	fmt.Println("  f - Fetch new results")
	fmt.Println("  e - Edit search parameters")
	fmt.Println("  d - Delete search")
	fmt.Println("  b - Back to search list")
	
	fmt.Print("\nEnter choice: ")
	scanner.Scan()
	choice := strings.TrimSpace(scanner.Text())
	
	fmt.Printf("[DEBUG] Search details - user selected: '%s'\n", choice)
	
	switch choice {
	case "f":
		fmt.Println("\n[COMING SOON] Fetch new results feature will be implemented soon.")
		fmt.Println("This will allow you to retrieve the latest articles from arXiv based on your search criteria.")
		pressEnterToContinue(scanner)
		return // Return to search list after action
	case "e":
		fmt.Println("\n[COMING SOON] Edit search parameters feature will be implemented soon.")
		fmt.Println("This will allow you to modify the search description, arXiv URL, and other parameters.")
		pressEnterToContinue(scanner)
		return // Return to search list after action
	case "d":
		deleteSearch(ctx, queries, searchID, searchStats.Description, scanner)
		return // Return to search list after action
	case "b":
		return // Return to search list
	default:
		fmt.Println("Invalid selection. Please choose from the options listed.")
		time.Sleep(1 * time.Second) // Brief pause
		return // Return to search list after invalid selection
	}
}

// defaultIfNullInt64 returns the default value if the SQL nullable value is not valid
func defaultIfNullInt64(value sql.NullInt64, defaultValue int64) int64 {
	if value.Valid {
		return value.Int64
	}
	return defaultValue
}

// deleteSearch handles the deletion of a search and its related records
func deleteSearch(ctx context.Context, queries *database.Queries, searchID interface{}, description string, scanner *bufio.Scanner) {
	fmt.Printf("[DEBUG] Starting deleteSearch function for ID: %v\n", searchID)
	
	// Show warning and confirmation
	fmt.Println("\n=== Delete Search ===")
	fmt.Printf("You are about to delete the search: \"%s\"\n", description)
	fmt.Println("This will permanently remove the search and all associated articles and candidates.")
	fmt.Print("Are you sure you want to continue? (y/n): ")
	
	scanner.Scan()
	confirm := strings.ToLower(strings.TrimSpace(scanner.Text()))
	
	if confirm != "y" && confirm != "yes" {
		fmt.Println("Deletion cancelled.")
		return
	}
	
	// Note: In a production environment, this should ideally be done in a transaction
	// to ensure atomicity, but we'll handle the deletes sequentially for simplicity
	
	// Step 1: First delete candidate associations (if we had the query)
	fmt.Println("Removing candidate associations...")
	// Note: We don't have the specific query yet, so we're just showing the message
	// In a complete implementation, you would call a query like:
	// err := queries.DeleteCandidateSearchesBySearchID(ctx, searchID)
	// if err != nil {
	//     fmt.Printf("Error removing candidate associations: %v\n", err)
	//     return
	// }
	
	// Step 2: Delete articles related to this search (if we had the query)
	fmt.Println("Removing related articles...")
	// Note: We don't have the specific query yet, so we're just showing the message
	// In a complete implementation, you would call a query like:
	// err := queries.DeleteArticlesBySearchID(ctx, searchID)
	// if err != nil {
	//     fmt.Printf("Error removing articles: %v\n", err)
	//     return
	// }
	
	// Step 3: Finally delete the search itself
	fmt.Println("Removing search record...")
	err := queries.DeleteSearch(ctx, searchID)
	if err != nil {
		fmt.Printf("Error deleting search: %v\n", err)
		return
	}
	
	fmt.Println("Search and all related data deleted successfully.")
}

// createNewSearch handles the creation of a new search
func createNewSearch(ctx context.Context, queries *database.Queries, scanner *bufio.Scanner) {
	fmt.Println("[DEBUG] Starting createNewSearch function")
	// Get search description
	fmt.Print("Enter search description: ")
	scanner.Scan()
	description := strings.TrimSpace(scanner.Text())
	if description == "" {
		fmt.Println("Search description cannot be empty.")
		return
	}

	// Get arXiv URL
	fmt.Print("Enter arXiv URL (e.g., http://rss.arxiv.org/rss/cs.LG+cs.PL): ")
	scanner.Scan()
	arxivURL := strings.TrimSpace(scanner.Text())
	if arxivURL == "" {
		fmt.Println("arXiv URL cannot be empty.")
		return
	}

	// Set up search parameters
	now := time.Now()
	defaultResultsPerFetch := sql.NullInt64{Int64: 50, Valid: true}

	// Create search record
	fmt.Println("[DEBUG] Creating search parameters")
	searchParams := database.CreateSearchParams{
		ID:              uuid.New(),
		CreatedAt:       now,
		UpdatedAt:       now,
		Description:     description,
		ArvixUrl:        arxivURL,
		ResultsPerFetch: defaultResultsPerFetch,
		LastFetchDate:   sql.NullTime{Valid: false},
	}
	fmt.Printf("[DEBUG] Search parameters created: ID=%v, Description=%s, URL=%s\n", 
		searchParams.ID, searchParams.Description, searchParams.ArvixUrl)

	fmt.Println("[DEBUG] Executing CreateSearch query...")
	search, err := queries.CreateSearch(ctx, searchParams)
	if err != nil {
		fmt.Printf("Error creating search: %v\n", err)
		return
	}

	fmt.Printf("Search created successfully with ID: %v\n", search.ID)
}
