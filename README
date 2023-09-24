# Portugal commute
I recently moved to Lisbon and was having a hard time finding apartments in cities that were in a radius from my job by train/subway.
Since I couldn't find any feature in Google Maps that allowed me to do exactly what I wanted (if someone knows about a feature that renders this program useless, please, tell me!), I coded this. At first in Python, but I got some [issues]() with Python's Vercel runtime and finally in TypeScript.
I welcome any tips, suggestions, PRs or issues 😊.

## Usage
Simply call https://portugal-commute.vercel.app/api/rail with the `location` (try an address. Google Maps will search it for you!), the `maxMinutes` you're willing to accept in your commute and your Google Maps `apiKey` (with permissions to use the Routes API) as query params and you gonna get a response with all the _concelhos_ that are in that travel distance radius.
The program assumes a search with the desired arrival time as 09:00h of the next day, just to be as accurate as possible.

## TODO
- Accept a list of address to have a more accurate search to the user's needs.
- Allow a custom `arrivalTime` as query params.
