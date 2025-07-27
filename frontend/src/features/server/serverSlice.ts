import { createAsyncThunk, createSlice } from '@reduxjs/toolkit'
import { RootState } from '../../app/store';
import { FetchStatus } from '../../app/types';
import { toast } from 'react-toastify';
import { changeCertbotStatusApi, editServerApi, getServerApi, getServerDetailsApi } from './serverApi';
import { ChangeSettingPayload, Server, ServerDetails, ServerFetchPayload, ServerSavePayload, ServerSaveRequest, ServerSettings } from './types';
import { CERTBOT_STATUS_SETTING } from './constants';

export interface ServerState {
    server: Server | null;
    serverDetails: ServerDetails | null;
    serverStatus: FetchStatus;
    serverDetailsStatus: FetchStatus;
    serverSaveStatus: FetchStatus;
    serverSettings: ServerSettings | null;
    changeCertbotStatusStatus: FetchStatus;
}

const initialState: ServerState = {
    server: null,
    serverDetails: null,
    serverStatus: FetchStatus.Idle,
    serverDetailsStatus: FetchStatus.Idle,
    serverSaveStatus: FetchStatus.Idle,
    serverSettings: null,
    changeCertbotStatusStatus: FetchStatus.Idle,
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

export const changeCertbotStatus = createAsyncThunk(
    'certbotstatus/change',
    async (payload: ChangeSettingPayload) => {
        const request = {
            guid: payload.guid,
            name: CERTBOT_STATUS_SETTING,
            value: payload.value,
            token: payload.token,
        };

        return await changeCertbotStatusApi(request);
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

                const settings = serverDetails?.settings;

                state.serverSettings = {
                    certbotStatus: settings && settings[CERTBOT_STATUS_SETTING] === 'true',
                };
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
            }).addCase(changeCertbotStatus.pending, state => {
                state.changeCertbotStatusStatus = FetchStatus.Pending;
            }).addCase(changeCertbotStatus.fulfilled, (state, action) => {
                state.changeCertbotStatusStatus = FetchStatus.Succeeded;

                if (state.serverSettings) {
                    state.serverSettings.certbotStatus = !state.serverSettings.certbotStatus;
                }

                if (action.payload) {
                    toast.success(`Certbot integration enabled. Certbot version: ${action.payload}`)
                } else {
                    toast.success('Certbot integration disabled')
                }
            }).addCase(changeCertbotStatus.rejected, (state, action) => {
                state.changeCertbotStatusStatus = FetchStatus.Failed;

                if (action.error.message) {
                    toast.error(action.error.message);
                }
            });
    },
});

export const selectServer = (state: RootState) => state.server.server;
export const selectServerFetchStatus = (state: RootState) => state.server.serverStatus;
export const selectServerDetails = (state: RootState) => state.server.serverDetails;
export const selectServerSettings = (state: RootState) => state.server.serverSettings;
export const selectServerDetailsFetchStatus = (state: RootState) => state.server.serverDetailsStatus;
export const selectServerSaveStatus = (state: RootState) => state.server.serverSaveStatus;
export const selectChangeCertbotStatusStatus = (state: RootState) => state.server.changeCertbotStatusStatus;

export default serverSlice.reducer
