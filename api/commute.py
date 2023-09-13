from blacksheep import Application
from blacksheep import Response
from blacksheep import json
from blacksheep import bad_request
from blacksheep import JSONContent
from blacksheep.client import ClientSession


ROUTES_LIMITING = 49


async def configure_http_client(app):
    http_client = ClientSession()
    app.services.add_instance(http_client)


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
) -> Response:
    reachable_destinations = []

    request: dict = {
        "origins": [],
        "destinations": [{"waypoint": {"address": location}}],
        "travelMode": "TRANSIT",
        "units": "METRIC",
        "transitPreferences": {
            "allowedTravelModes": ["RAIL"],
        },
    }

    municipalities: list[str] = await get_municipalities(http_client=http_client)

    for i in range(0, len(municipalities), ROUTES_LIMITING):
        slice = municipalities[i : i + ROUTES_LIMITING]

        request["origins"] = [{"waypoint": {"address": m}} for m in slice]
        print(request)

        url = f"https://routes.googleapis.com/distanceMatrix/v2:computeRouteMatrix?key={api_key}&fields=duration,distanceMeters,originIndex"

        response = await http_client.post(url, content=JSONContent(request))
        assert response is not None

        data = await response.json()
        print(data)

        for result in data:
            if "duration" not in result:
                continue
            commute_time_seconds = result["duration"][:-1]
            commute_time_minutes = int(commute_time_seconds) // 60

            if commute_time_minutes <= max_minutes:
                original_request = request["origins"][result["originIndex"]]

                reachable_destinations.append(
                    {
                        "concelho": original_request["waypoint"]["address"],
                        "duration_in_minutes": commute_time_minutes,
                    }
                )

    if not reachable_destinations:
        return bad_request("No routes found. Try a different search query.")

    return json(reachable_destinations)
