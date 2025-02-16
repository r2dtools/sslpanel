import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import { RootState } from './store';
import { ColorTheme } from '../types/theme';

export interface AppState {
    colorMode: ColorTheme;
}

const initialState: AppState = {
    colorMode: ColorTheme.Dark,
}

export const appSlice = createSlice({
    name: 'app',
    initialState,
    reducers: {
        setAppColorMode: (state, action: PayloadAction<ColorTheme>) => {
            state.colorMode = action.payload;
        },
    },
});

export const { setAppColorMode } = appSlice.actions;

export const selectColorMode = (state: RootState) => state.app.colorMode;

export default appSlice.reducer
