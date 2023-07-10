#!/usr/bin/env python3

import argparse
import requests
import sys
from bs4 import BeautifulSoup

parser = argparse.ArgumentParser(description='Scrape donorbox page for donation goal status')
parser.add_argument(
    '--interval',
    '-i',
    help='Scrape interval',
    type=int,
    default=60,
    required=False
)
args = parser.parse_args()

url = 'https://donorbox.org/support-black-girls-code/fundraiser/christopher-dunaj'
headers = {"User-Agent": "Just someone learning python"}

try:
    page = requests.get(url, headers=headers)
except:
    print('There was a problem reaching the URL: {0}'.format(url))
    sys.exit(1)

try:
    page.raise_for_status()
except Exception as exc:
    print('There was a problem with the response from the URL: {0}'.format(exc))
    sys.exit(2)

soup = BeautifulSoup(page.content, 'html.parser')

try:
    progress = str(soup.find(attrs={"id": "panel-1"}).contents[1])
    num_donators = soup.find(attrs={"id": "paid-count"}).text
    #'0'
    total_raised = soup.find(attrs={"id": "total-raised"}).text
    #'$0'
    raise_goal = soup.find(attrs={"id": "panel-1"}).contents[1].find_all(attrs={"class": "bold"})[2].text
    #'$500'
except:
    print('Did not find expected contents at that URL.')
    sys.exit(3)

if progress:
    #print('Progress found!')
    #print(progress)
    print('Number of donators: {0}'.format(num_donators))
    print('Total raised: {0}'.format(total_raised))
    print('Funraiser goal: {0}'.format(raise_goal))
