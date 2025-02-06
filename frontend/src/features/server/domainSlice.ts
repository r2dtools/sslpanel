import { createAsyncThunk, createSlice, PayloadAction } from '@reduxjs/toolkit';
import { FetchStatus } from '../../app/types';
import {
    CommonChallengeDirStatusChangePayload,
    Domain,
    DomainFetchPayload,
    DomainSecurePayload,
    DomainSettings,
    DomainSettingsPayload,
} from './types';
import { RootState } from '../../app/store';
import { getServerDomainApi } from './serverApi';
import { toast } from 'react-toastify';
import { changeCommonDirStatusApi, getCommonDirApi, secureDomainApi } from './domainApi';

export interface DomainState {
    domain: Domain | null;
    settings: DomainSettings | null;
    domainStatus: FetchStatus;
    domainSecureStatus: FetchStatus;
    settingsStatus: FetchStatus;
    changeCommonDirStatusStatus: FetchStatus;
}

const initialState: DomainState = {
    domain: null,
    settings: null,
    domainStatus: FetchStatus.Idle,
    domainSecureStatus: FetchStatus.Idle,
    settingsStatus: FetchStatus.Idle,
    changeCommonDirStatusStatus: FetchStatus.Idle,
}

export const fetchServerDomain = createAsyncThunk(
    'domain',
    async (payload: DomainFetchPayload) => {
        return await getServerDomainApi(payload.guid, payload.domainName, payload.token);
    },
);

export const secureServerDomain = createAsyncThunk(
    'secure',
    async (payload: DomainSecurePayload) => {
        const request = {
            email: payload.email,
            subjects: payload.subjects,
            servername: payload.servername,
            webserver: payload.webserver,
            docroot: payload.docroot,
            challengetype: payload.challengetype,
            assign: payload.assign,
        };

        await secureDomainApi(payload.guid, request, payload.token);

        return await getServerDomainApi(payload.guid, payload.servername, payload.token);
    },
);

export const fetchSettings = createAsyncThunk(
    'settings',
    async (payload: DomainSettingsPayload) => {
        const commonDirRequest = {
            servername: payload.domain.servername,
            webserver: payload.domain.webserver,
        };

        const status = await getCommonDirApi(payload.guid, commonDirRequest, payload.token);

        return {
            commondirstatus: status,
            renewal: false,
        };
    },
);

export const changeCommonDirStatus = createAsyncThunk(
    'commondirstatus/change',
    async (payload: CommonChallengeDirStatusChangePayload) => {
        const request = {
            servername: payload.domain.servername,
            webserver: payload.domain.webserver,
            status: payload.status,
        };

        await changeCommonDirStatusApi(payload.guid, request, payload.token);

        return true;
    },
);

export const domainSlice = createSlice({
    name: 'domain',
    initialState,
    reducers: {
        domainFetched: (state, action: PayloadAction<Domain>) => {
            state.domainStatus = FetchStatus.Succeeded;
            state.domain = action.payload;
        },
    },
    extraReducers: builder => {
        builder
            .addCase(fetchServerDomain.pending, state => {
                state.domainStatus = FetchStatus.Pending;
            })
            .addCase(fetchServerDomain.fulfilled, (state, action) => {
                state.domainStatus = FetchStatus.Succeeded;
                state.domain = action.payload;
            })
            .addCase(fetchServerDomain.rejected, (state, action) => {
                state.domainStatus = FetchStatus.Failed;
                state.domain = null;

                if (action.error.message) {
                    toast.error(action.error.message);
                }
            }).addCase(secureServerDomain.pending, state => {
                state.domainSecureStatus = FetchStatus.Pending;
            })
            .addCase(secureServerDomain.fulfilled, (state, action) => {
                state.domainSecureStatus = FetchStatus.Succeeded;
                state.domain = action.payload;
            })
            .addCase(secureServerDomain.rejected, (state, action) => {
                state.domainSecureStatus = FetchStatus.Failed;

                if (action.error.message) {
                    toast.error(action.error.message);
                }
            }).addCase(fetchSettings.pending, state => {
                state.settingsStatus = FetchStatus.Pending;
            })
            .addCase(fetchSettings.fulfilled, (state, action) => {
                state.settingsStatus = FetchStatus.Succeeded;
                state.settings = action.payload;
            })
            .addCase(fetchSettings.rejected, (state, action) => {
                state.settingsStatus = FetchStatus.Failed;
                state.settings = null;

                if (action.error.message) {
                    toast.error(action.error.message);
                }
            }).addCase(changeCommonDirStatus.pending, state => {
                state.changeCommonDirStatusStatus = FetchStatus.Pending;
            })
            .addCase(changeCommonDirStatus.fulfilled, state => {
                state.changeCommonDirStatusStatus = FetchStatus.Succeeded;

                if (state.settings) {
                    state.settings.commondirstatus = {
                        status: !state.settings.commondirstatus.status,
                    };
                }
            }).addCase(changeCommonDirStatus.rejected, (state, action) => {
                state.changeCommonDirStatusStatus = FetchStatus.Failed;

                if (action.error.message) {
                    toast.error(action.error.message);
                }
            });
    },
});

export const { domainFetched } = domainSlice.actions;

export const selectDomain = (state: RootState) => state.domain.domain;
export const selectSettings = (state: RootState) => state.domain.settings;
export const selectDomainFetchStatus = (state: RootState) => state.domain.domainStatus;
export const selectSettingsFetchStatus = (state: RootState) => state.domain.settingsStatus;
export const selectDomainSecureStatus = (state: RootState) => state.domain.domainSecureStatus;
export const selectChangeCommonDirStatusStatus = (state: RootState) => state.domain.changeCommonDirStatusStatus;

export default domainSlice.reducer
