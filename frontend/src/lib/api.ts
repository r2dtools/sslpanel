import axios, { AxiosError, AxiosRequestConfig } from 'axios';

const api = axios.create({
    baseURL: import.meta.env.VITE_API_HOST,
    headers: {
        'Content-Type': 'application/json',
    },
});

interface AxiosErrorResponseData {
    code: number;
    message: string;
};

export const getErrorMessage = (error: unknown) => {
    const defaultErrorMessage = 'Unknown error. Please try again later';

    if (error instanceof AxiosError && error.response) {
        const data: AxiosErrorResponseData | '' = error.response.data;

        return data === '' ? defaultErrorMessage : data.message;
    }

    if (error instanceof Error) {
        return error.message || defaultErrorMessage;
    }

    return String(error || defaultErrorMessage);
}

export const configWithAuth = (token: string, config?: AxiosRequestConfig): AxiosRequestConfig => {
    config = config || {};
    const headers = config.headers || {};

    return {
        ...config,
        headers: {
            ...headers,
            Authorization: `Bearer ${token}`,
        },
    };
};

api.interceptors.response.use((response) => response, (error) => {
    if (error.response.status === 401) {
        (window as Window).location = '/auth/sigin';
    }

    return Promise.reject(error);
});

export default api;
