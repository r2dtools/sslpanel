import type { Action, ThunkAction } from '@reduxjs/toolkit';
import { configureStore } from '@reduxjs/toolkit';
import authReducer from '../features/auth/authSlice';
import serversReducer from '../features/server/serversSlice';
import serverReducer from '../features/server/serverSlice';
import domainReducer from '../features/domain/domainSlice';

export const store = configureStore({
    reducer: {
        auth: authReducer,
        servers: serversReducer,
        server: serverReducer,
        domain: domainReducer,
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
