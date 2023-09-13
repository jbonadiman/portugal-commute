from blacksheep import Application
from blacksheep.client import ClientSession
from blacksheep import JSONContent
from typing import Union


async def configure_http_client(app):
    http_client = ClientSession()
    app.services.add_instance(http_client)  # register a singleton


async def dispose_http_client(app):
    http_client = app.service_provider.get(ClientSession)
    await http_client.close()


app = Application()
app.on_start += configure_http_client
app.on_stop += dispose_http_client
get = app.router.get


municipalities_cache: list[str] = []


async def get_municipalities(http_client: ClientSession) -> list[str]:
    global municipalities_cache

    if not municipalities_cache:
        response = await http_client.get("https://json.geoapi.pt/municipios")

        assert response is not None
        municipalities_cache = await response.json()

    return municipalities_cache


@get("/")
async def index(
    http_client: ClientSession, location: str, max_minutes: int, api_key: str
) -> list[dict[str, Union[str, int]]]:
    reachable_destinations = []

    request: dict[str, Union[dict[str, Union[str, list[str]]], str]] = {
        "origin": {"address": ""},
        "destination": {"address": location},
        "travelMode": "TRANSIT",
        "units": "METRIC",
        "transitPreferences": {
            "allowedTravelModes": ["RAIL"],
        },
    }

    for municipality in await get_municipalities(http_client=http_client):
        print(municipality)

        request["origin"] = {"address": municipality}

        url = f"https://routes.googleapis.com/directions/v2:computeRoutes?key={api_key}&fields=routes.duration,routes.distanceMeters"
        response = await http_client.post(url, content=JSONContent(request))
        assert response is not None

        data = await response.json()

        if "routes" not in data:
            continue

        commute_time_seconds = data["routes"][0]["duration"][:-1]
        print(commute_time_seconds)
        commute_time_minutes = int(commute_time_seconds) // 60
        print(commute_time_minutes)

        if commute_time_minutes <= max_minutes:
            reachable_destinations.append(
                {"concelho": municipality, "duration": commute_time_minutes}
            )

    return reachable_destinations
