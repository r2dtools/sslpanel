import { createAsyncThunk, createSlice, PayloadAction } from '@reduxjs/toolkit';
import { FetchStatus } from '../../app/types';
import { Domain, DomainFetchPayload } from './types';
import { RootState } from '../../app/store';
import { getServerDomainApi } from './serverApi';
import { toast } from 'react-toastify';

export interface DomainState {
    domain: Domain | null;
    domainStatus: FetchStatus;
    domainSecureStatus: FetchStatus;
}

const initialState: DomainState = {
    domain: null,
    domainStatus: FetchStatus.Idle,
    domainSecureStatus: FetchStatus.Idle,
}

export const fetchServerDomain = createAsyncThunk(
    'domain',
    async (payload: DomainFetchPayload) => {
        return await getServerDomainApi(payload.guid, payload.domainName, payload.token);
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
        builder.addCase(fetchServerDomain.pending, state => {
            state.domainStatus = FetchStatus.Pending;
        })
            .addCase(fetchServerDomain.fulfilled, (state, action) => {
                state.domainStatus = FetchStatus.Succeeded;
                state.domain = action.payload;
            })
            .addCase(fetchServerDomain.rejected, (state, action) => {
                state.domainSecureStatus = FetchStatus.Failed;
                state.domain = null;

                if (action.error.message) {
                    toast.error(action.error.message);
                }
            });
    },
});

export const { domainFetched } = domainSlice.actions;

export const selectDomain = (state: RootState) => state.domain.domain;
export const selectDomainFetchStatus = (state: RootState) => state.domain.domainStatus;
export const selectDomainSecureStatus = (state: RootState) => state.domain.domainSecureStatus;

export default domainSlice.reducer
