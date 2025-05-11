export interface ChangePasswordRequest {
    password: string;
    newPassword: string;
    token: string;
};

export interface ChangePasswordPayload {
    password: string;
    newPassword: string;
    token: string;
};

export enum ColorTheme {
    Light = 'light',
    Dark = 'dark',
};
