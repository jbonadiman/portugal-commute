import type { VercelResponse } from '@vercel/node';

const ROUTES_LIMIT = 49;

let municipalitiesCache: string[] = [];

type Route = {
    duration: string;
    distance: number;
    originIndex: number;
}

type RequestQuery = {
    location: string;
    maxMinutes: number;
    apiKey: string;
}

type MapsRequest = {
    origins: { waypoint: { address: string } }[];
    destinations: { waypoint: { address: string } }[];
    travelMode: string;
    units: string;
    transitPreferences: { allowedTravelModes: string[]; };
}

type Result = {
    concelho: string;
    durationMinutes: number;
}

const getMunicipalities = async () => {
    if (!municipalitiesCache.length) {
        const res = await fetch("https://json.geoapi.pt/municipios");
        municipalitiesCache = await res.json();
    }

    return municipalitiesCache;
}

export default async function (request: { query: RequestQuery }, response: VercelResponse) {
    const {
        location,
        maxMinutes,
        apiKey
    } = request.query

    if (!location || !maxMinutes || !apiKey) {
        return response.status(400).json({ error: 'Missing required parameters' });
    }

    const maxSeconds = maxMinutes * 60;

    const gMapsEndpoint = `https://routes.googleapis.com/distanceMatrix/v2:computeRouteMatrix?key=${apiKey}&fields=duration,distanceMeters,originIndex`

    const mapsRequest: MapsRequest = {
        origins: [],
        destinations: [{ waypoint: { address: location } }],
        travelMode: 'TRANSIT',
        units: 'METRIC',
        transitPreferences: { allowedTravelModes: ['RAIL'] }
    }

    const municipalities = await getMunicipalities();
    const reachableDestinations: Result[] = [];

    for (let i = 0; i < municipalities.length; i += ROUTES_LIMIT) {
        const slice = municipalities.slice(i, i + ROUTES_LIMIT);

        mapsRequest.origins = slice.map(m => ({ waypoint: { address: m } }));

        const res = await fetch(gMapsEndpoint, {
            method: 'POST',
            body: JSON.stringify(mapsRequest)
        })

        const data = await res.json();

        reachableDestinations.push(...data
            .filter((route: Route) => route.duration)
            .map((route: Route) => {
                route.duration = route.duration.replace('s', '');
                return route;
            })
            .filter((route: Route) => (Number(route.duration)) <= maxSeconds)
            .map((route: Route) => {
                return {
                    concelho: municipalities[i + route.originIndex],
                    durationMinutes: Math.round(Number(route.duration) / 60)
                };
            }));
        console.log(JSON.stringify(reachableDestinations))
    }

    return response.json(reachableDestinations);
}
