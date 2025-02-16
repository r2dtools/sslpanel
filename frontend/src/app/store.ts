import type { Action, ThunkAction } from '@reduxjs/toolkit';
import { configureStore } from '@reduxjs/toolkit';
import appReducer from './appSlice';
import authReducer from '../features/auth/authSlice';
import serversReducer from '../features/server/serversSlice';
import serverReducer from '../features/server/serverSlice';
import domainReducer from '../features/domain/domainSlice';
import certificatesReducer from '../features/certificate/certificatesSlice';

export const store = configureStore({
    reducer: {
        app: appReducer,
        auth: authReducer,
        servers: serversReducer,
        server: serverReducer,
        domain: domainReducer,
        certificates: certificatesReducer,
    }
})

export type AppStore = typeof store
export type RootState = ReturnType<AppStore['getState']>

export type AppDispatch = AppStore['dispatch']

export type AppThunk<ThunkReturnType = void> = ThunkAction<
    ThunkReturnType,
    RootState,
    unknown,
    Action
>
