import React from 'react';
import ReactDOM from 'react-dom/client';
import { BrowserRouter as Router } from 'react-router-dom';
import { Provider } from 'react-redux';
import App from './App';
import { store } from './app/store';
import './css/style.css';
import './css/satoshi.css';
import { createTheme, Flowbite } from 'flowbite-react';

const customTheme = createTheme({
    textInput: {
        field: {
            input: {
                colors: {
                    gray: 'border-gray-300 bg-gray-50 text-gray-900 placeholder-gray-500 focus:border-blue-500 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-700 dark:text-white dark:placeholder-gray-400 dark:focus:border-blue-500 dark:focus:ring-blue-500',
                },
            },
        }
    },
    tabs: {
        tablist: {
            tabitem: {
                base: 'flex items-center justify-center rounded-t-lg p-4 text-sm font-medium first:ml-0 focus:outline-none disabled:cursor-not-allowed disabled:text-gray-400 disabled:dark:text-gray-500',
                variant: {
                    underline: {
                        active: {
                            'on': 'rounded-t-lg border-b-2 border-blue-600 text-blue-600 dark:border-blue-500 dark:text-blue-500',
                        },
                    },
                },
            }
        },
    }
});

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
    <React.StrictMode>
        <Provider store={store}>
            <Router>
                <Flowbite theme={{ theme: customTheme }}>
                    <App />
                </Flowbite>
            </Router>
        </Provider>
    </React.StrictMode>,
);
