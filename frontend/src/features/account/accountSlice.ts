import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import { FetchStatus } from '../../app/types';
import { changePasswordApi } from './accountApi';
import { ChangePasswordPayload } from './types';
import { toast } from 'react-toastify';
import { RootState } from '../../app/store';

export interface AccountState {
    passwordChangeStatus: FetchStatus;
};

const initialState = {
    passwordChangeStatus: FetchStatus.Idle,
};

export const changePassword = createAsyncThunk(
    'account/change-assword',
    async (payload: ChangePasswordPayload) => {
        return await changePasswordApi({
            password: payload.password,
            newPassword: payload.newPassword,
            token: payload.token,
        });
    },
);

export const accountSlice = createSlice({
    name: 'account',
    initialState,
    reducers: {},
    extraReducers: builder => {
        builder
            .addCase(changePassword.pending, state => {
                state.passwordChangeStatus = FetchStatus.Pending;
            })
            .addCase(changePassword.fulfilled, state => {
                state.passwordChangeStatus = FetchStatus.Succeeded;
                toast.success('Password changed successfully');
            })
            .addCase(changePassword.rejected, (state, action) => {
                state.passwordChangeStatus = FetchStatus.Failed;

                if (action.error.message) {
                    toast.error(action.error.message);
                }
            });
    },
});

export const selectPasswordChangeStatus = (state: RootState) => state.account.passwordChangeStatus;

export default accountSlice.reducer;
