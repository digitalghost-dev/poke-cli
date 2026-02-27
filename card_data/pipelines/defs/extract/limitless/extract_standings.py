import dagster as dg
import polars as pl
import requests
from bs4 import BeautifulSoup

def seasons(year: int) -> list[dict]:
    url = "https://labs.limitlesstcg.com/"
    r = requests.get(url)
    soup = BeautifulSoup(r.content, 'html.parser')

    season_header = soup.find("h2", string=lambda x: x and f"{year}" in x)
    tournament_list = season_header.find_next_sibling("ul")

    tournaments = []
    for a in tournament_list.find_all("a", href=True):
        parts = list(a.stripped_strings)
        tournaments.append({
            "name": parts[0],
            "date": parts[-1],
            "link": f"https://labs.limitlesstcg.com{a['href']}",
        })

    return tournaments


def build_standings(tournament: dict) -> pl.DataFrame | None:
    tournament_id = tournament["link"].split("/")[-2]

    r = requests.get(tournament["link"])
    soup = BeautifulSoup(r.content, 'html.parser')

    table = soup.find("table", class_="data-table striped")

    if table:
        headers = ['Rank', 'Name', 'Country', 'Points', 'Record', 'OPW%', 'OOPW%', 'Deck', 'Decklist', 'Unknown']

        rows = []
        tbody = table.find('tbody')

        for tr in tbody.find_all('tr'):
            cells = tr.find_all('td')

            if len(cells) == 1:
                continue

            row_data = []
            for i, td in enumerate(cells):
                if i == 2:
                    img = td.find('img')
                    if img:
                        country = img.get('alt') or img.get('title') or ''
                        row_data.append(country)
                    else:
                        row_data.append('')

                elif i == 7:  # Deck column
                    pokemon_imgs = td.find_all('img', class_='pokemon')
                    if pokemon_imgs:
                        pokemon_names = [img.get('alt', '') for img in pokemon_imgs if img.get('alt')]
                        pokemon_string = '/'.join(pokemon_names)
                        row_data.append(pokemon_string)
                    else:
                        row_data.append('')

                elif i == 8:  # Decklist column
                    link = td.find('a')
                    if link:
                        decklist_url = link.get('href', '')
                        row_data.append(f"https://labs.limitlesstcg.com{decklist_url}" if decklist_url else '')
                    else:
                        row_data.append('')

                else:
                    cell_text = td.get_text(strip=True)
                    row_data.append(cell_text)

            rows.append(row_data)

        df = pl.DataFrame(rows, schema=headers, orient="row")

        df = df.drop("Unknown")

        df = df.with_columns(pl.lit(tournament_id).alias("tournament_id"))

        df = df.rename(
            {
                "Rank": "rank",
                "Name": "name",
                "Country": "country",
                "Points": "points",
                "Record": "record",
                "OPW%": "opp_win_percent",
                "OOPW%": "opp_opp_win_percent",
                "Deck": "deck",
                "Decklist": "decklist",
            },
        )

        df = df.cast({"rank": pl.Int16, "points": pl.Int16})

        return df

    return None


@dg.asset(kinds={"Polars"}, name="create_standings_dataframe")
def create_standings_dataframe() -> pl.DataFrame:
    tournaments = seasons(2026)

    dfs = []
    for t in tournaments:
        df = build_standings(t)
        if df is not None:
            dfs.append(df)
            print(f"Loaded {df.shape[0]} rows from {t['name']}")

    return pl.concat(dfs)
