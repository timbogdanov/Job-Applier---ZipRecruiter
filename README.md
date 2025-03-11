# Job Applier – ZipRecruiter

This Go-based automation script streamlines applying to jobs on [ZipRecruiter](https://www.ziprecruiter.com) by simulating human-like interactions. It navigates through job listings, highlights key elements, clicks the "1-Click Apply" button where available, handles pop-up modals, and moves on to the next listing automatically.

## Features

- **Automated Navigation**  
  Iterates through job cards on ZipRecruiter, highlighting each processed card.

- **1-Click Apply Automation**  
  Locates the "1-Click Apply" button, applies with a force click if needed, and waits for any resulting modals to appear.

- **Modal Handling**  
  Closes additional information or confirmation pop-ups automatically to ensure seamless progression.

- **Human-Like Behavior**  
  Uses random delays and slight mouse movements to mimic human interaction and reduce the chance of detection.

## Prerequisites

1. **Go (1.16 or newer)**
2. **Google Chrome**  
   Must be launched with a remote debugging port.
3. **Playwright for Go**  
   Community fork: [github.com/playwright-community/playwright-go](https://github.com/playwright-community/playwright-go)

## Installation

1. **Clone the Repository**  
```bash
git clone git@github.com:timbogdanov/Job-Applier---ZipRecruiter.git
cd Job-Applier---ZipRecruiter
```

2. **Install Dependencies**  
Make sure your Go module is initialized, then install missing packages:
```bash
go mod tidy
```

3. **Install Playwright Browsers (Optional)**  
If you haven’t installed them globally, you may need:
```bash
npx playwright install
```

## Usage

1. **Launch Chrome with Remote Debugging**  
    Make sure Chrome is closed, then run:
```bash
chrome --remote-debugging-port=9222
```
    Replace `chrome` with the appropriate path to your Chrome executable if needed.

2. **Log into ZipRecruiter**  
    In that same Chrome instance, go to [ZipRecruiter](https://www.ziprecruiter.com) and manually log in.

3. **Adjust Script Settings (Optional)**  
   - Open `main.go` and modify:
     - The **search URL** or job search criteria.
     - **Sleep intervals** if you want different random delays.
     - **CSS highlight styles**, etc.

4. **Run the Script**  
```bash
go run main.go
```

   The script connects to your Chrome instance on port `9222` and begins processing job cards automatically.

## How It Works

1. **Navigates to the Provided URL**  
   The script visits the ZipRecruiter job search page specified in `main.go`.

2. **Finds Job Cards**  
   It highlights each job card, clicks it, and opens details in the right pane.

3. **Checks for “1-Click Apply”**  
   If the button is found, the script highlights it, clicks it, then waits briefly.

4. **Closes Pop-Ups**  
   If there is an info modal requesting additional data, or a "cancel application" confirmation, the script automatically closes them.

5. **Pagination**  
   Once all visible job cards are processed, it clicks the "Next Page" link and continues until there are no more pages.

## Troubleshooting

- **Not Connecting to Chrome**  
  Verify Chrome is running with `--remote-debugging-port=9222`.
- **Login Issues**  
  Make sure you’re logged into ZipRecruiter in that same instance of Chrome.
- **Selector Changes**  
  ZipRecruiter’s site may update. If elements can’t be found, adjust the selectors in `main.go`.
- **Timeouts**  
  Increase timeouts if the script fails due to slow network or site loading.

## Contributing

Pull requests and issue reports are welcome! If you have suggestions, bug fixes, or new ideas, feel free to open an issue or submit a PR.

## License

This project is licensed under the **MIT License**.