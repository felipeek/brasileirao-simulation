#!/usr/bin/env python3
import subprocess
import re
import json
import os
from collections import Counter

def get_all_teams(directory):
    teams = set()
    for filename in os.listdir(directory):
        if filename.endswith(".json"):
            with open(os.path.join(directory, filename), 'r') as f:
                data = json.load(f)
                if 'Name' in data:
                    teams.add(data['Name'])
    return teams

def run_process_and_get_champion():
    # Prepare the command
    command = ['./bin/bsim']
    command.append('--non-interactive')
    
    # Run the process and capture the output
    result = subprocess.run(command, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
    output = result.stdout
    
    # Find the champion using regex
    match = re.search(r'The champion: \[([^\]]+)\]!', output)
    if match:
        return match.group(1)
    else:
        return None

def main():
    # Get all teams from the JSON files
    all_teams = get_all_teams('./teams')

    champions = []

    # Run the process 1000 times
    for i in range(1000):
        champion = run_process_and_get_champion()
        if champion:
            champions.append(champion)
        else:
            print(f"Run {i+1}: No champion found")

    # Count occurrences of each champion
    champion_count = Counter(champions)

    # Ensure all teams are included in the final ranking
    for team in all_teams:
        if team not in champion_count:
            champion_count[team] = 0

    # Print the ranking
    print("Ranking of champions from most to least:")
    for team, count in champion_count.most_common():
        print(f"{team}: {count} times")

if __name__ == "__main__":
    main()

