
const GOOGLE_MAPS_LINK = (lat, lon) => {
    return `https://www.google.com/maps/search/?api=1&query=${lat},${lon}`;
};

export function getLocationData(locations = []) {
    let arr = [];
    locations.forEach(location => {
        arr.push({
            lat: location.lat,
            lon: location.lon,
            link: GOOGLE_MAPS_LINK(location.lat, location.lon),
            city: location.address.city || location.address.state,
            country: location.address.country
        });
    });

    return arr;
}