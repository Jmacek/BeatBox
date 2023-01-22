import React, { useState } from "react";
import TextField from '@mui/material/TextField';

import {apiUrl} from "../../configs/constants";
import { Song } from "../../configs/structures";
import SongList from "../component/SongList";

const SearchPage = () => {
    const [searchTerm, setSearchTerm] = useState("");
    const [songs, setSongs] = useState<Song[]>([]);

    const fetchData = async () => {
        const requestHeaders: HeadersInit = new Headers();
        requestHeaders.set('Content-Type', 'application/json');

        const response = await fetch(`${apiUrl}/api/music/search?searchTerm=${searchTerm}`, {
            headers: requestHeaders,
        });
        const data = await response.json();
        setSongs(data);
    }

    const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setSearchTerm(event.target.value);
    };

    return (
      <>
        <TextField
          label="Search - Full Titles Only (case sensitive)"
          variant="outlined"
          fullWidth
          value={searchTerm}
          onChange={handleChange}
          onKeyDown={event => {
            if (event.key === 'Enter') {
                fetchData();
            }
          }}
        />
      <SongList songs={songs} />
    </>
    );
};

export default SearchPage;
