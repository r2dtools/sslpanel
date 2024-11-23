import { createAsyncThunk, createSlice } from '@reduxjs/toolkit'
import { RootState } from '../../app/store';
import { FetchStatus } from '../../app/types';
import { toast } from 'react-toastify';
import { addServerApi, deleteServerApi, editServerApi, getServersApi } from './serverApi';
import { ServerSavePayload, Server, SaveServerRequest, ServerDeletePayload } from './types';

export interface ServerState {
    servers: Server[];
    serversStatus: FetchStatus;
    serverSaveStatus: FetchStatus;
    serverDeleteStatus: FetchStatus;
}

const initialState: ServerState = {
    servers: [],
    serversStatus: FetchStatus.Idle,
    serverSaveStatus: FetchStatus.Idle,
    serverDeleteStatus: FetchStatus.Idle,
}

export const fetchServers = createAsyncThunk(
    'server/list',
    async (token: string) => {
        return await getServersApi(token);
    },
);

export const addServer = createAsyncThunk(
    'server/add',
    async (payload: ServerSavePayload) => {
        const request: SaveServerRequest = {
            name: payload.name,
            ipv4_address: payload.ipv4_address,
            ipv6_address: payload.ipv6_address,
            agent_port: payload.agent_port,
        };

        await addServerApi(request, payload.token);

        return await getServersApi(payload.token);
    },
);

export const editServer = createAsyncThunk(
    'server/edit',
    async (payload: ServerSavePayload) => {
        const request: SaveServerRequest = {
            name: payload.name,
            ipv4_address: payload.ipv4_address,
            ipv6_address: payload.ipv6_address,
            agent_port: payload.agent_port,
        };

        await editServerApi(payload.id as number, request, payload.token);

        return await getServersApi(payload.token);
    },
);

export const deleteServer = createAsyncThunk(
    'server/delete',
    async (payload: ServerDeletePayload) => {
        await deleteServerApi(payload.id, payload.token);

        return await getServersApi(payload.token);
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
            }).addCase(addServer.pending, state => {
                state.serverSaveStatus = FetchStatus.Pending;
            }).addCase(addServer.fulfilled, (state, action) => {
                state.serverSaveStatus = FetchStatus.Succeeded;
                state.servers = action.payload;
            }).addCase(addServer.rejected, (state, action) => {
                state.serverSaveStatus = FetchStatus.Failed;

                if (action.error.message) {
                    toast.error(action.error.message);
                }
            }).addCase(editServer.pending, state => {
                state.serverSaveStatus = FetchStatus.Pending;
            }).addCase(editServer.fulfilled, (state, action) => {
                state.serverSaveStatus = FetchStatus.Succeeded;
                state.servers = action.payload;
            }).addCase(editServer.rejected, (state, action) => {
                state.serverSaveStatus = FetchStatus.Failed;

                console.log(action);
                if (action.error.message) {
                    toast.error(action.error.message);
                }
            }).addCase(deleteServer.pending, state => {
                state.serverDeleteStatus = FetchStatus.Pending;
            }).addCase(deleteServer.fulfilled, (state, action) => {
                state.serverDeleteStatus = FetchStatus.Succeeded;
                state.servers = action.payload;
            }).addCase(deleteServer.rejected, (state, action) => {
                state.serverDeleteStatus = FetchStatus.Failed;

                if (action.error.message) {
                    toast.error(action.error.message);
                }
            });
    },
});

export const selectServers = (state: RootState) => state.server.servers;
export const selectServersFetchStatus = (state: RootState) => state.server.serversStatus;
export const selectServerSaveStatus = (state: RootState) => state.server.serverSaveStatus;

export default serverSlice.reducer
