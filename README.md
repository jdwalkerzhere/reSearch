# reSearch

reSearch is intended to initially be a CLI-based AI Research Talent exploration tool (I may eventually make a web interface), but for the time being I'm keeping things simple.

## Dependencies:
- SQLite: We persist the searches, articles seen, and candidates to try and avoid duplicate work
- Anthropic API: The program assumes you're being a responsible adult and storing this in a `.env` file

## Usage: 
The intended workflow is that when the user invokes the application is: 
- Options presentation ("Create New Search", "Manage Searches", "Check New Results", "Fetch Older Results", "Delete Search", etc)
- If creating a new search, have a brief conversation with an LLM to determine what arXiv categories could be relevant
    - For new searches, we start making requests to arXiv against the desired categories (default is from most recent to least recent)
    - For each research paper we pass the summary to an LLM to evaluate its relevance
    - If the paper is deemed relevant (getting this dailed in will likely be an iterative process), we attempt to find the LinkedIn or GitHub profiles (or both if available) and store them for exploration later.

## Current State:
- Non-functional: So far I've not been able to get web search tool use integrated into the Agent and hardcoding search for a POC seems out of scope/possibly more error prone.
