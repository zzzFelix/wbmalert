# wbm-alert
This script creates (text) snapshots for a given list of websites. A request to each website is made every 30 seconds. If the contents of the website have changed, a notification sound is played.

## Limitations
The snapshot only considers the textual contents of a website. If images, links, or HTML attributes change, the alarm may not sound. In more technical terms: `<a href="https://google.com">Link</a>` is sanitized to `Link`. All whitespace (including tabs and line breaks) is also removed.

## What is it good for?
I used it to monitor real estate websites to be notified of new listings. I was particularly interested in Berlin's state-owned Wohnungsbaugesellschaft Berlin-Mitte (WBM), hence the name `wbm-alert`.

## Prequisites
Go >= `1.21`, older versions may work but aren't tested.

## Get started
- Clone the repo
- Adjust `configuration.json`: Edit the `websites` array to contain your own links.
- Optionally adjust the interval
- Install and run the program `go install && wbmalert -c configuration.json`

## Configuration options
- `interval`: Time in seconds between requests
- `websites`: Array of websites to make requests to.
    - `name`: Name to identify the website. Does not need to be unique.
    - `url`: Url to make HTTP GET request to.
    - `regexRemove` (optional): A regular expression. Removes every substring that matches.