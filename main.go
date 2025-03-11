package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

func main() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Start Playwright
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not launch playwright: %v", err)
	}

	// Connect to an existing Chrome instance on port 9222
	browser, err := pw.Chromium.ConnectOverCDP("http://localhost:9222")
	if err != nil {
		log.Fatalf("could not connect to existing Chrome instance: %v\nEnsure Chrome is running with --remote-debugging-port=9222", err)
	}

	// Get the first context (assuming one tab is open)
	contexts := browser.Contexts()
	if len(contexts) == 0 {
		log.Fatalf("no browser contexts found; ensure Chrome is open with a tab")
	}
	context := contexts[0]

	// Get all pages (tabs)
	pages := context.Pages()
	if len(pages) == 0 {
		log.Fatalf("no pages found in the browser context")
	}

	// Find a suitable page (preferably on ziprecruiter.com)
	var page playwright.Page
	for _, p := range pages {
		url := p.URL()
		log.Printf("Checking tab with URL: %s", url)
		if strings.Contains(url, "ziprecruiter.com") {
			page = p
			log.Printf("Selected ZipRecruiter tab: %s", url)
			break
		}
		if !strings.HasPrefix(url, "chrome-extension://") && !strings.HasPrefix(url, "chrome://") {
			// Fallback to any non-extension, non-chrome page if no ZipRecruiter tab found yet
			page = p
			log.Printf("Selected fallback tab (not an extension): %s", url)
		}
	}

	if page == nil {
		log.Fatalf("No suitable tab found. Please open a tab on 'ziprecruiter.com' or any non-extension page and try again.")
	}

	// Log connection and current state
	log.Println("Connected to your existing Chrome window. Taking over immediately.")
	currentURL := page.URL()
	log.Printf("Current URL before navigation: %s", currentURL)

	// Check if logged in (optional safeguard)
	if !strings.Contains(currentURL, "ziprecruiter.com") || strings.Contains(currentURL, "authn/login") {
		log.Println("Warning: You might not be logged in. Please ensure you're on a ZipRecruiter page post-login.")
	}

	// --- Step 1: Go to the job search page ---
	jobSearchURL := "https://www.ziprecruiter.com/jobs-search?search=PHP+Developer&location=Remote+%28USA%29&days=5&refine_by_employment=employment_type%3Aall&refine_by_salary=&refine_by_salary_ceil=&lvk=UqnX-FDhizyUPjp9N7KSjg.--Nk1TOdf2c"
	log.Printf("Navigating to: %s", jobSearchURL)
	_, err = page.Goto(jobSearchURL, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
		Timeout:   playwright.Float(30000), // 30-second timeout
	})
	if err != nil {
		log.Printf("Navigation failed: %v", err)
		log.Println("Current page content (for debugging):")
		content, _ := page.Content()
		log.Println(content[:500]) // Print first 500 chars of HTML
		log.Fatalf("Stopping due to navigation error.")
	}
	log.Println("Successfully navigated to job search page.")
	randomSleep(3*time.Second, 7*time.Second)

	// --- Step 2: Process job cards and apply ---
	for {
		simulateHumanBehavior(page) // Add human-like behavior on each page
		jobCards := page.Locator(".job_results_two_pane .job_result_two_pane")
		count, err := jobCards.Count()
		if err != nil {
			log.Fatalf("could not count job cards: %v", err)
		}

		for i := 0; i < count; i++ {
			log.Printf("Processing card #%d of %d", i+1, count)

			// Highlight (style) job card by setting its background color
			cardLocator := jobCards.Nth(i)
			highlightElement(cardLocator, "background-color: #D3FFD9; border: 1px solid #39FF50;")
			elementHandle, err := cardLocator.ElementHandle()
			if err != nil || elementHandle == nil {
				log.Printf("failed to get element handle for card #%d: %v", i+1, err)
				continue
			}
			if err = simulateClick(page, elementHandle); err != nil {
				log.Printf("failed to click card #%d: %v", i+1, err)
				continue
			}
			randomSleep(1*time.Second, 3*time.Second)

			// Highlight right pane: set background and border styles
			rightPane := page.Locator("[data-testid='right-pane']")
			highlightElement(rightPane, "background-color: #D3FFD9; border: 1px solid #39FF50;")
			if err := rightPane.WaitFor(playwright.LocatorWaitForOptions{
				State:   playwright.WaitForSelectorStateVisible,
				Timeout: playwright.Float(5000),
			}); err != nil {
				log.Printf("Right pane not visible for card #%d: %v", i+1, err)
				continue
			}

			boundingBox, err := rightPane.BoundingBox()
			if err != nil {
				log.Printf("Failed to get bounding box for right pane on card #%d: %v", i+1, err)
				continue
			}

			// Calculate the center coordinates of the right pane
			x := boundingBox.X + boundingBox.Width/2
			y := boundingBox.Y + boundingBox.Height/2

			mouse := page.Mouse()
			// Move the mouse to the center of the right pane
			if err := mouse.Move(x, y); err != nil {
				log.Printf("Failed to move mouse for card #%d: %v", i+1, err)
				continue
			}
			time.Sleep(1 * time.Second)

			// Highlight target button and update its style for "1-Click Apply"
			targetButton := rightPane.GetByRole("button").Filter(playwright.LocatorFilterOptions{
				HasText: "1-Click Apply",
			})
			log.Printf("Looking for '1-Click Apply' button")

			// Click with force and timeout
			err = targetButton.Click(playwright.LocatorClickOptions{
				Force:   playwright.Bool(true),
				Timeout: playwright.Float(2000),
			})
			if err != nil {
				log.Printf("Failed to click button: %v", err)
			} else {
				log.Printf("Successfully clicked '1-Click Apply' button")
				// Wait a couple of seconds after clicking to allow any modal to appear
				time.Sleep(2 * time.Second)

				// Look for the first modal (the info modal) with a close button
				infoModalCloseButton := page.Locator("div[role='dialog'] button[aria-label='Close']").First()
				btnCount, err := infoModalCloseButton.Count()
				if err == nil && btnCount > 0 {
					log.Printf("Info modal detected. Attempting to close it.")
					err = infoModalCloseButton.Click(playwright.LocatorClickOptions{
						Force:   playwright.Bool(true),
						Timeout: playwright.Float(5000),
					})
					if err != nil {
						log.Printf("Failed to click info modal close button: %v", err)
					} else {
						log.Printf("Closed the info modal.")
					}
					// Wait briefly to allow any subsequent modal to appear
					time.Sleep(1 * time.Second)

					// Check for the confirmation modal that appears after closing the info modal
					confirmModal := page.Locator("div[aria-label='Are you sure you want to cancel?']")
					// We use IsVisible() on the locator (it returns a bool) so we check if it's visible
					if visible, _ := confirmModal.IsVisible(); visible {
						cancelButton := confirmModal.Locator("button:has-text('Cancel Application')").First()
						cancelCount, err := cancelButton.Count()
						if err == nil && cancelCount > 0 {
							log.Printf("Confirmation modal detected. Clicking 'Cancel Application' to close it.")
							err = cancelButton.Click(playwright.LocatorClickOptions{
								Force:   playwright.Bool(true),
								Timeout: playwright.Float(5000),
							})
							if err != nil {
								log.Printf("Failed to click 'Cancel Application' button: %v", err)
							} else {
								log.Printf("Closed the confirmation modal.")
							}
						} else {
							log.Printf("Confirmation modal did not contain a 'Cancel Application' button.")
						}
					} else {
						log.Printf("No confirmation modal detected after closing the info modal.")
					}
				} else {
					log.Printf("No info modal appeared after clicking '1-Click Apply'.")
				}
				// Wait a couple of seconds before moving to the next job card
				time.Sleep(2 * time.Second)
			}

			randomSleep(3*time.Second, 6*time.Second)
			log.Println("----------------------------------------------------------------------------------------------")
		}

		// Highlight next page link
		nextPage := page.Locator("a[title='Next Page']")
		highlightElement(nextPage, "background-color: #D3FFD9; border: 1px solid #39FF50;")
		if nextPage == nil {
			log.Println("No next page found. Exiting loop.")
			break
		}
		if err = nextPage.Click(); err != nil {
			log.Fatalf("could not click next page button: %v", err)
		}
		if err = page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
			State: playwright.LoadStateNetworkidle,
		}); err != nil {
			log.Fatalf("pagination page load failed: %v", err)
		}
		randomSleep(3*time.Second, 7*time.Second)
	}

	log.Println("Finished processing all pages.")

	// Clean up
	if err = browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}
	if err = pw.Stop(); err != nil {
		log.Fatalf("could not stop playwright: %v", err)
	}
}

// highlightElement applies CSS to an element via JavaScript
func highlightElement(locator playwright.Locator, css string) {
	_, err := locator.Evaluate("el => el.style.cssText += '"+css+"'", nil)
	if err != nil {
		log.Printf("Failed to apply highlight CSS: %v", err)
	}
}

// randomSleep introduces a random delay between min and max
func randomSleep(min, max time.Duration) {
	time.Sleep(min + time.Duration(rand.Int63n(int64(max-min))))
}

// simulateClick mimics a human click with mouse movement
func simulateClick(page playwright.Page, element playwright.ElementHandle) error {
	// Highlight the element being clicked by setting its border (for visual feedback)
	_, err := element.Evaluate("el => el.style.cssText += 'border: 1px solid #39FF50;'", nil)
	if err != nil {
		log.Printf("Failed to highlight element for click: %v", err)
	}

	box, err := element.BoundingBox()
	if err != nil {
		return err
	}
	if box == nil {
		return fmt.Errorf("element not visible")
	}

	x := box.X + box.Width/2 + float64(rand.Intn(10)-5)
	y := box.Y + box.Height/2 + float64(rand.Intn(10)-5)
	if err := page.Mouse().Move(x, y); err != nil {
		return err
	}
	randomSleep(100*time.Millisecond, 300*time.Millisecond)

	// Click near the top of the element to avoid clicking company name
	y = box.Y + (box.Height * 0.2) + float64(rand.Intn(10)-5)
	if err := page.Mouse().Click(x, y); err != nil {
		return err
	}
	return nil
}

// simulateHumanBehavior adds scrolling and random mouse movements
func simulateHumanBehavior(page playwright.Page) {
	// Random scroll (no element to highlight here, but could highlight scrollable area if desired)
	page.Evaluate(fmt.Sprintf(`window.scrollBy(0, %d);`, rand.Intn(500)-250))
	randomSleep(500*time.Millisecond, 1500*time.Millisecond)

	// Random mouse movement across the page (no specific element to highlight)
	x := float64(rand.Intn(800) + 200) // Between 200 and 1000
	y := float64(rand.Intn(400) + 100) // Between 100 and 500
	page.Mouse().Move(x, y)
	randomSleep(300*time.Millisecond, 800*time.Millisecond)
}
