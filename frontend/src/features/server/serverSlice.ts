import { createAsyncThunk, createSlice } from '@reduxjs/toolkit'
import { RootState } from '../../app/store';
import { FetchStatus } from '../../app/types';
import { toast } from 'react-toastify';
import { editServerApi, getServerApi, getServerDetailsApi } from './serverApi';
import { Server, ServerDetails, ServerFetchPayload, ServerSavePayload, ServerSaveRequest } from './types';

export interface ServerState {
    server: Server | null;
    serverDetails: ServerDetails | null;
    serverStatus: FetchStatus;
    serverDetailsStatus: FetchStatus;
    serverSaveStatus: FetchStatus;
}

const initialState: ServerState = {
    server: null,
    serverDetails: null,
    serverStatus: FetchStatus.Idle,
    serverDetailsStatus: FetchStatus.Idle,
    serverSaveStatus: FetchStatus.Idle,
}

export const fetchServer = createAsyncThunk(
    'server/info',
    async (payload: ServerFetchPayload) => {
        return await getServerApi(payload.guid, payload.token);
    },
);

export const fetchServerDetails = createAsyncThunk(
    'server/details',
    async (payload: ServerFetchPayload) => {
        return await getServerDetailsApi(payload.guid, payload.token);
    },
);

export const editServer = createAsyncThunk(
    'server/edit',
    async (payload: ServerSavePayload) => {
        const request: ServerSaveRequest = {
            name: payload.name,
            ipv4_address: payload.ipv4_address,
            ipv6_address: payload.ipv6_address,
            agent_port: payload.agent_port,
            token: payload.token,
        };

        await editServerApi(payload.id as number, request, payload.authToken);

        return await getServerApi(payload.guid as string, payload.authToken);
    },
);

export const serverSlice = createSlice({
    name: 'server',
    initialState,
    reducers: {},
    extraReducers: builder => {
        builder
            .addCase(fetchServer.pending, state => {
                state.serverStatus = FetchStatus.Pending;
            })
            .addCase(fetchServer.fulfilled, (state, action) => {
                state.serverStatus = FetchStatus.Succeeded;
                state.server = action.payload;
            })
            .addCase(fetchServer.rejected, (state, action) => {
                state.serverStatus = FetchStatus.Failed;
                state.server = null;

                if (action.error.message) {
                    toast.error(action.error.message);
                }
            }).addCase(fetchServerDetails.pending, state => {
                state.serverDetailsStatus = FetchStatus.Pending;
            })
            .addCase(fetchServerDetails.fulfilled, (state, action) => {
                state.serverDetailsStatus = FetchStatus.Succeeded;

                const serverDetails = action.payload;
                serverDetails.domains = serverDetails.domains.sort((a, b) => a.servername.localeCompare(b.servername));
                state.serverDetails = serverDetails;
            })
            .addCase(fetchServerDetails.rejected, (state, action) => {
                state.serverDetailsStatus = FetchStatus.Failed;
                state.serverDetails = null;

                if (action.error.message) {
                    toast.error(action.error.message);
                }
            }).addCase(editServer.pending, state => {
                state.serverSaveStatus = FetchStatus.Pending;
            }).addCase(editServer.fulfilled, (state, action) => {
                state.serverSaveStatus = FetchStatus.Succeeded;
                state.server = action.payload;
            }).addCase(editServer.rejected, (state, action) => {
                state.serverSaveStatus = FetchStatus.Failed;

                if (action.error.message) {
                    toast.error(action.error.message);
                }
            });
    },
});

export const selectServer = (state: RootState) => state.server.server;
export const selectServerFetchStatus = (state: RootState) => state.server.serverStatus;
export const selectServerDetails = (state: RootState) => state.server.serverDetails;
export const selectServerDetailsFetchStatus = (state: RootState) => state.server.serverDetailsStatus;
export const selectServerSaveStatus = (state: RootState) => state.server.serverSaveStatus;

export default serverSlice.reducer
