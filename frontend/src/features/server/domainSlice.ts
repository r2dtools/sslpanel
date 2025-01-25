import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { FetchStatus } from '../../app/types';
import { Domain } from './types';
import { RootState } from '../../app/store';

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

export const domainSlice = createSlice({
    name: 'domain',
    initialState,
    reducers: {
        domainFetched: (state, action: PayloadAction<Domain>) => {
            state.domain = action.payload;
            state.domainStatus = FetchStatus.Succeeded;
        },
    },
});

export const { domainFetched } = domainSlice.actions;

export const selectDomain = (state: RootState) => state.domain.domain;
export const selectDomainFetchStatus = (state: RootState) => state.domain.domainStatus;
export const selectDomainSecureStatus = (state: RootState) => state.domain.domainSecureStatus;

export default domainSlice.reducer
