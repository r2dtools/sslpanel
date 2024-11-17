import { createAsyncThunk, createSlice } from '@reduxjs/toolkit'
import { RootState } from '../../app/store';
import { FetchStatus } from '../../app/types';
import { toast } from 'react-toastify';
import { getServers } from './serverApi';
import { Server } from './types';

export interface ServerState {
    servers: Server[];
    serversStatus: FetchStatus;
}

const initialState: ServerState = {
    servers: [],
    serversStatus: FetchStatus.Idle,
}

export const fetchServers = createAsyncThunk(
    'server/list',
    async (token: string) => {
        return await getServers(token);
    },
);

export const serverSlice = createSlice({
    name: 'server',
    initialState,
    reducers: {},
    extraReducers: builder => {
        builder
            .addCase(fetchServers.pending, state => {
                state.serversStatus = FetchStatus.Pending;
            })
            .addCase(fetchServers.fulfilled, (state, action) => {
                state.serversStatus = FetchStatus.Succeeded;
                state.servers = action.payload;
            })
            .addCase(fetchServers.rejected, (state, action) => {
                state.serversStatus = FetchStatus.Failed;

                if (action.error.message) {
                    toast.error(action.error.message);
                }
            });
    },
});

export const selectServers = (state: RootState) => state.server.servers;
export const selectServersFetchStatus = (state: RootState) => state.server.serversStatus;

export default serverSlice.reducer
