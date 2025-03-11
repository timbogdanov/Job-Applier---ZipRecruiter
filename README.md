ZipRecruiter Job Apply Automation

This automation script, written in Go with Playwright, automatically processes job listings on ZipRecruiter by applying using the "1-Click Apply" feature. It highlights elements to provide visual feedback, handles potential modal pop-ups gracefully, and simulates human-like interactions to enhance stability.

Prerequisites

Go (version 1.16 or higher)

Google Chrome browser

Playwright Go library

Installation

Clone the repository:

git clone https://github.com/your-username/ziprecruiter-automation.git
cd ziprecruiter-automation

Install Dependencies:

go mod tidy

Install Playwright Browsers (if you haven't already):

npx playwright install

Prepare Chrome for Automation

Start Chrome with remote debugging enabled:

chrome --remote-debugging-port=9222

Manually navigate to ZipRecruiter and log into your account.

Before Running the Script

Open Chrome with remote debugging on port 9222:

chrome --remote-debugging-port=9222

Navigate manually to ZipRecruiter and log in.

Verify you're logged in and are able to browse job listings.

Running the Automation Script

Once logged in and the Chrome instance is active, run:

go run main.go

Workflow

The script navigates through each job card, highlights each processed job card in #39FF50.

The right pane showing job details is highlighted with a background color of #D3FFD9 and a border of #39FF50.

It searches for the "1-Click Apply" button and highlights it in #39FF50.

Upon clicking the button, it waits briefly and automatically closes any modals that request additional info, allowing continuous unattended processing.

Handling Modals

If a modal requesting more information appears after clicking apply, the script automatically closes it.

If a confirmation modal appears afterward (e.g., "Are you sure you want to cancel?"), it automatically confirms cancellation and continues.

Running the Script

Execute the automation by running:

go run main.go

Configuration and Adjustments

To change the job search URL or the search criteria, update the URL inside main.go.

Adjust random sleep durations within the script to simulate different human-like interaction timings if necessary.

Troubleshooting

Remote debugging not working: Ensure Chrome is correctly started with --remote-debugging-port=9222.

Not logged in: Make sure you're logged into ZipRecruiter before running the script.

Selectors not found: Verify that the selectors used in the script still match ZipRecruiter’s current HTML structure.

Contributions

Feel free to open issues or submit pull requests if you have improvements or find issues with the script.

License

MIT License.