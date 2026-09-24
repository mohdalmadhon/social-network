
// this function fetch the location data from api 
//
// @returns array of objects
export async function getSuggestedLocations(searchValue = "") {
    const resp = await fetch(
        `/api/location/search?q=${encodeURIComponent(searchValue)}`,
        {
            method: "GET"
        }
    );

    if (!resp.ok) {
        throw new Error("could not fetch locations");
    }

    const data = await resp.json();


    return data;
}

export function getLocationData(results = []) {
    if (!Array.isArray(results)) return []

    return results
        .map((item) => {
            const address = item.address || {}

            return {
                city: address.city || address.town || address.village || address.municipality || '',
                state: address.state || address.region || '',
                country: address.country || '',
                lat: item.lat,
                lon: item.lon,
                link: String(item.place_id || `${item.lat}-${item.lon}`)
            }
        })
        .filter((location) => (location.city || location.state) && location.country)
}