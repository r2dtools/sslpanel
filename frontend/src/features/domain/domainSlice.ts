import { createAsyncThunk, createSlice, PayloadAction } from '@reduxjs/toolkit';
import { FetchStatus } from '../../app/types';
import {
    ChangeSettingPayload,
    Domain,
    DomainConfigPayload,
    DomainFetchPayload,
    DomainSecurePayload,
    DomainSettings,
    DomainSettingsPayload,
} from './types';
import { RootState } from '../../app/store';
import { toast } from 'react-toastify';
import {
    changeCommonDirStatusApi,
    getCommonDirApi,
    getDomainConfigApi,
    getDomainApi,
    secureDomainApi,
    getDomainSettingsApi,
    changeSettingApi,
} from './domainApi';
import { COMMON_DIR_SETTING, RENEWAL_SETTING } from './constants';

export interface DomainState {
    domain: Domain | null;
    settings: DomainSettings | null;
    config: string | null;
    domainStatus: FetchStatus;
    domainSecureStatus: FetchStatus;
    settingsStatus: FetchStatus;
    changeCommonDirStatusStatus: FetchStatus;
    changeRenewalStatus: FetchStatus;
    configStatus: FetchStatus;
}

const initialState: DomainState = {
    domain: null,
    settings: null,
    config: null,
    domainStatus: FetchStatus.Idle,
    domainSecureStatus: FetchStatus.Idle,
    settingsStatus: FetchStatus.Idle,
    changeCommonDirStatusStatus: FetchStatus.Idle,
    changeRenewalStatus: FetchStatus.Idle,
    configStatus: FetchStatus.Idle,
}

export const fetchServerDomain = createAsyncThunk(
    'domain',
    async (payload: DomainFetchPayload) => {
        const request = {
            guid: payload.guid,
            domainname: payload.domainname,
            token: payload.token
        };
        return await getDomainApi(request);
    },
);

export const secureServerDomain = createAsyncThunk(
    'secure',
    async (payload: DomainSecurePayload) => {
        const request = {
            guid: payload.guid,
            email: payload.email,
            subjects: payload.subjects,
            servername: payload.servername,
            webserver: payload.webserver,
            challengetype: payload.challengetype,
            assign: payload.assign,
            token: payload.token,
        };

        await secureDomainApi(request);

        return await getDomainApi({
            guid: payload.guid,
            domainname: payload.servername,
            webserver: payload.webserver,
            token: payload.token,
        });
    },
);

export const fetchSettings = createAsyncThunk(
    'settings',
    async (payload: DomainSettingsPayload) => {
        const request = {
            guid: payload.guid,
            domainname: payload.domain.servername,
            webserver: payload.domain.webserver,
            token: payload.token,
        };

        const status = await getCommonDirApi(request);
        const settings = await getDomainSettingsApi(request);

        return {
            commondirstatus: status,
            renewalstatus: settings[RENEWAL_SETTING] === 'true',
        };
    },
);

export const changeCommonDirStatus = createAsyncThunk(
    'commondirstatus/change',
    async (payload: ChangeSettingPayload) => {
        const request = {
            guid: payload.guid,
            servername: payload.domain.servername,
            webserver: payload.domain.webserver,
            name: COMMON_DIR_SETTING,
            status: payload.status,
            token: payload.token,
        };

        await changeCommonDirStatusApi(request);
    },
);

export const changeRenewal = createAsyncThunk(
    'renewal/change',
    async (payload: ChangeSettingPayload) => {
        const request = {
            guid: payload.guid,
            servername: payload.domain.servername,
            webserver: payload.domain.webserver,
            name: RENEWAL_SETTING,
            status: payload.status,
            token: payload.token,
        };

        await changeSettingApi(request);
    },
);

export const fetchConfig = createAsyncThunk(
    'config',
    async (payload: DomainConfigPayload) => {
        const request = {
            guid: payload.guid,
            webserver: payload.domain.webserver,
            domainname: payload.domain.servername,
            token: payload.token,
        };

        return await getDomainConfigApi(request);
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
        configReset: state => {
            state.configStatus = FetchStatus.Idle;
            state.config = null;
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
            }).addCase(fetchConfig.pending, state => {
                state.configStatus = FetchStatus.Pending;
                state.config = null;
            })
            .addCase(fetchConfig.fulfilled, (state, action) => {
                state.configStatus = FetchStatus.Succeeded;
                state.config = action.payload;
            })
            .addCase(fetchConfig.rejected, (state, action) => {
                state.configStatus = FetchStatus.Failed;
                state.config = null;

                if (action.error.message) {
                    toast.error(action.error.message);
                }
            }).addCase(changeRenewal.pending, state => {
                state.changeRenewalStatus = FetchStatus.Pending;
            })
            .addCase(changeRenewal.fulfilled, state => {
                state.changeRenewalStatus = FetchStatus.Succeeded;

                if (state.settings) {
                    state.settings.renewalstatus = !state.settings.renewalstatus;
                }
            }).addCase(changeRenewal.rejected, (state, action) => {
                state.changeRenewalStatus = FetchStatus.Failed;

                if (action.error.message) {
                    toast.error(action.error.message);
                }
            });
    },
});

export const { domainFetched, configReset } = domainSlice.actions;

export const selectDomain = (state: RootState) => state.domain.domain;
export const selectSettings = (state: RootState) => state.domain.settings;
export const selectConfig = (state: RootState) => state.domain.config;
export const selectDomainFetchStatus = (state: RootState) => state.domain.domainStatus;
export const selectSettingsFetchStatus = (state: RootState) => state.domain.settingsStatus;
export const selectDomainSecureStatus = (state: RootState) => state.domain.domainSecureStatus;
export const selectChangeCommonDirStatusStatus = (state: RootState) => state.domain.changeCommonDirStatusStatus;
export const selectConfigFetchStatus = (state: RootState) => state.domain.configStatus;
export const selectChangeRenewalStatus = (state: RootState) => state.domain.changeRenewalStatus;

export default domainSlice.reducer
