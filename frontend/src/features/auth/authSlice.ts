import { createAsyncThunk, createSlice } from '@reduxjs/toolkit'
import { User } from './types'
import { me } from './authApi';
import { RootState } from '../../app/store';
import { FetchStatus } from '../../app/types';
import { toast } from 'react-toastify';

export interface AuthState {
    currentUser: User | null;
    status: FetchStatus;
}

const initialState: AuthState = {
    currentUser: null,
    status: FetchStatus.Idle,
}

export const fetchCurrentUser = createAsyncThunk(
    'auth/me',
    async (token: string) => {
        return await me(token);
    },
);

export const authSlice = createSlice({
    name: 'auth',
    initialState,
    reducers: {
        logout: state => {
            state.currentUser = null;
        },
    },
    extraReducers: builder => {
        builder
            .addCase(fetchCurrentUser.pending, state => {
                state.status = FetchStatus.Pending;
            })
            .addCase(fetchCurrentUser.fulfilled, (state, action) => {
                state.status = FetchStatus.Succeeded;
                state.currentUser = action.payload;
            })
            .addCase(fetchCurrentUser.rejected, (state, action) => {
                state.status = FetchStatus.Failed;

                if (action.error.message) {
                    toast.error(action.error.message);
                }
            });
    },
});

export const { logout } = authSlice.actions;

export const selectCurrentUser = (state: RootState) => state.auth.currentUser;
export const selectCurrentUserFetchStatus = (state: RootState) => state.auth.status;

export default authSlice.reducer
