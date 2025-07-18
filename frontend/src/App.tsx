import React, { useEffect } from 'react';
import { Route, Routes, useLocation, useNavigate } from 'react-router-dom';
import Loader from './components/Loader/GlobalLoader';
import PageTitle from './components/PageTitle';
import DefaultLayout from './layout/DefaultLayout';
import AuthContext from './features/auth/context';
import SignIn from "./pages/Authentication/SignIn";
import SignUp from "./pages/Authentication/SignUp";
import Settings from "./pages/Settings";
import useLocalStorage from './hooks/useLocalStorage';
import { ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import { useAppDispatch, useAppSelector } from './app/hooks';
import { fetchCurrentUser, selectCurrentUser, selectCurrentUserFetchStatus } from './features/auth/authSlice';
import useColorMode from './hooks/useColorMode';
import { FetchStatus, RouteItem } from './app/types';
import AuthLayout from './layout/AuthLayout';
import ServerList from './pages/ServerList/ServerList';
import Server from './pages/Server/Server';
import useAuthToken from './features/auth/hooks';
import RoutesContext from './app/context';
import Error404 from './pages/Error404';
import Domain from './pages/Domain/Domain';
import CertificateList from './pages/CertificateList/CertificatesList';
import DomainList from './pages/DomainList/DomainList';
import { setAppColorMode } from './app/appSlice';
import { ColorTheme } from './features/account/types';
import Confirm from './pages/Authentication/Confirm';
import Recover from './pages/Authentication/Recover';
import Reset from './pages/Authentication/Reset';
import InstallSslBot from './pages/Help/InstallSslBot';

function App() {
    const { pathname } = useLocation();
    const navigate = useNavigate();
    const [colorMode] = useColorMode();

    const currentUser = useAppSelector(selectCurrentUser);
    const currentUserLoadStatus = useAppSelector(selectCurrentUserFetchStatus);
    const [authToken, setAuthToken] = useAuthToken();
    const [authTokenExpire, setAuthTokenExpire] = useLocalStorage<string | null>("r2panel-token-expire", null);
    const dispatch = useAppDispatch();

    useEffect(() => {
        dispatch(setAppColorMode(colorMode as ColorTheme));
    }, []);

    useEffect(() => {
        window.scrollTo(0, 0);
    }, [pathname]);

    useEffect(() => {
        if (!authToken || currentUserLoadStatus === FetchStatus.Pending) {
            return;
        }

        if (authTokenExpire) {
            const expire = Date.parse(authTokenExpire);

            if (expire <= Date.now()) {
                handleSignOut();

                return;
            }
        }

        dispatch(fetchCurrentUser(authToken));
    }, [authToken, authTokenExpire]);

    const handleSignIn = (token: string, expire: string) => {
        setAuthToken(token);
        setAuthTokenExpire(expire);
    };

    const handleSignOut = () => {
        setAuthToken(null);
        setAuthTokenExpire(null);
    }

    const getPageTitle = (title: string) => `${title} | SSL Control Panel`;

    const routes: { layout: React.FC<any>, props?: React.ComponentProps<any>, routes: RouteItem[] }[] = [
        {
            layout: DefaultLayout,
            props: {
                onSignOut: handleSignOut,
            },
            routes: [
                {
                    path: '*',
                    title: 'Page not found',
                    component: <Error404 />
                },
                {
                    path: "/settings",
                    title: getPageTitle('Settings'),
                    component: <Settings />
                },
                {
                    path: "/servers",
                    title: getPageTitle('Servers'),
                    name: "Servers",
                    component: <ServerList />
                },
                {
                    path: "/servers/:guid",
                    title: getPageTitle('Server'),
                    name: 'Server',
                    component: <Server />,
                },
                {
                    path: "/servers/:guid/domains",
                    title: getPageTitle('Domains'),
                    name: 'Domains',
                    component: <DomainList />,
                },
                {
                    path: "/servers/:guid/certificates",
                    title: getPageTitle('Certificates'),
                    name: 'Certificates',
                    component: <CertificateList />,
                },
                {
                    path: "/servers/:guid/domains/:name",
                    title: getPageTitle('Domain'),
                    name: 'Domain',
                    component: <Domain />,
                },
                {
                    path: "/documentation/install-sslbot",
                    title: getPageTitle('Install SSLBot '),
                    name: 'Install SSLBot',
                    component: <InstallSslBot />,
                },
            ],
        },
        {
            layout: AuthLayout,
            props: {},
            routes: [
                {
                    path: "/auth/signin",
                    public: true,
                    title: getPageTitle('Sign In'),
                    component: <SignIn onSignIn={handleSignIn} />
                },
                {
                    path: "/auth/signup",
                    public: true,
                    title: getPageTitle('Sign Up'),
                    component: <SignUp />
                },
                {
                    path: "/auth/confirm",
                    public: true,
                    title: getPageTitle('Email confirmation'),
                    component: <Confirm />
                },
                {
                    path: "/auth/recover",
                    public: true,
                    title: getPageTitle('Password recovery'),
                    component: <Recover />
                },
                {
                    path: "/auth/reset",
                    public: true,
                    title: getPageTitle('Password reset'),
                    component: <Reset />
                },
            ],
        },
    ];

    let allRoutes: RouteItem[] = [];
    routes.forEach(({ routes }) => allRoutes = allRoutes.concat(routes));

    const authRoutes = ["/auth/signin", "/auth/signup", "/auth/confirm"];

    useEffect(() => {
        const currentRoute = allRoutes.find((item: RouteItem) => item.path === pathname);

        if (
            !authToken
            && currentUserLoadStatus !== FetchStatus.Pending
            && currentUser === null
            && !currentRoute?.public
        ) {
            navigate("/auth/signin");

            return;
        }

        const isAuthRoute = authRoutes.includes(currentRoute?.path || "");

        if ((currentUser !== null && isAuthRoute) || currentRoute === null || pathname === '/') {
            navigate("/servers");
        }
    }, [currentUser, routes, authRoutes, pathname]);

    const filterRoutes = (item: RouteItem) => (currentUser !== null && !authRoutes.includes(item.path || "")) || (currentUser === null && item.public);

    const loading = currentUserLoadStatus === FetchStatus.Pending
        || (currentUser !== null && authRoutes.includes(pathname))
        || pathname === '/';

    return loading ? (
        <Loader />
    ) : (
        <AuthContext.Provider value={currentUser}>
            <RoutesContext.Provider value={allRoutes.filter(filterRoutes)}>
                <Routes>
                    {
                        routes.map(({ layout: Layout, routes, props }) => (
                            <Route element={<Layout {...props} />} key={Layout.name}>
                                {
                                    routes
                                        .filter(filterRoutes)
                                        .map(({ path, index, title, component }: RouteItem) => (
                                            <Route
                                                key={path}
                                                index={index}
                                                path={path}
                                                element={
                                                    <>
                                                        <PageTitle title={title} />
                                                        {component}
                                                    </>
                                                }
                                            />
                                        ))
                                }
                            </Route>
                        ))
                    }
                </Routes>
                <ToastContainer theme={colorMode === ColorTheme.Dark ? 'dark' : 'light'} />
            </RoutesContext.Provider>
        </AuthContext.Provider>
    );
}

export default App;
