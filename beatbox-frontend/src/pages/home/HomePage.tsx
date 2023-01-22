import React, { useState, useEffect } from 'react';

import { apiUrl } from "../../configs/constants";
import { Song } from "../../configs/structures";

import SongList from "../component/SongList";

const HomePage = () => {
    const [songs, setSongs] = useState<Song[]>([]);

    const fetchData = async () => {
        const requestHeaders: HeadersInit = new Headers();
        requestHeaders.set('Content-Type', 'application/json');

        const response = await fetch(apiUrl + '/api/music/songs', {
            headers: requestHeaders,
        });
        const data = await response.json();
        setSongs(data);
    }

    useEffect(() => {
        fetchData();
    }, []);

    return (
        <SongList songs={songs} />)
    ;
};

export default HomePage;