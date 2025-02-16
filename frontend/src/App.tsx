import React, { useEffect } from 'react';
import { Route, Routes, useLocation, useNavigate } from 'react-router-dom';
import Loader from './components/Loader/GlobalLoader';
import PageTitle from './components/PageTitle';
import DefaultLayout from './layout/DefaultLayout';
import AuthContext from './features/auth/context';
import SignIn from "./pages/Authentication/SignIn";
import SignUp from "./pages/Authentication/SignUp";
import Calendar from "./pages/Calendar";
import Chart from "./pages/Chart";
import ECommerce from "./pages/Dashboard/ECommerce";
import FormElements from "./pages/Form/FormElements";
import FormLayout from "./pages/Form/FormLayout";
import Profile from "./pages/Profile";
import Settings from "./pages/Settings";
import Tables from "./pages/Tables";
import Alerts from "./pages/UiElements/Alerts";
import Buttons from "./pages/UiElements/Buttons";
import useLocalStorage from './hooks/useLocalStorage';
import { ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import { useAppDispatch, useAppSelector } from './app/hooks';
import { fetchCurrentUser, selectCurrentUser, selectCurrentUserFetchStatus } from './features/auth/authSlice';
import useColorMode from './hooks/useColorMode';
import { FetchStatus, RouteItem } from './app/types';
import { ColorTheme } from './types/theme';
import AuthLayout from './layout/AuthLayout';
import ServerList from './pages/ServerList/ServerList';
import Server from './pages/Server/Server';
import useAuthToken from './features/auth/hooks';
import RoutesContext from './app/context';
import Error404 from './pages/Error404';
import Domain from './pages/Domain/Domain';
import CertificateList from './pages/CertificateList/CertificatesList';
import DomainList from './pages/DomainList/DomainList';

function App() {
    const { pathname } = useLocation();
    const navigate = useNavigate();
    const [colorMode] = useColorMode();

    const currentUser = useAppSelector(selectCurrentUser);
    const currentUserLoadStatus = useAppSelector(selectCurrentUserFetchStatus);
    const [authToken, setAuthToken] = useAuthToken();
    const [authTokenExpire, setAuthTokenExpire] = useLocalStorage<string | null>("r2panel-token-expire", null);
    const dispatch = useAppDispatch();

    const currentUserLoading = currentUserLoadStatus === FetchStatus.Pending;
    const loading = currentUserLoading;

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
                    path: '/auth/*',
                    title: "eCommerce Dashboard | TailAdmin - Tailwind CSS Admin Dashboard Template",
                    component: <ECommerce />
                },
                {
                    index: true,
                    title: "eCommerce Dashboard | TailAdmin - Tailwind CSS Admin Dashboard Template",
                    component: <ECommerce />
                },
                {
                    path: "/calendar",
                    title: "Calendar | TailAdmin - Tailwind CSS Admin Dashboard Template",
                    component: <Calendar />
                },
                {
                    path: "/profile",
                    title: "Profile | TailAdmin - Tailwind CSS Admin Dashboard Template",
                    component: <Profile />
                },
                {
                    path: "/forms/form-elements",
                    title: "Form Elements | TailAdmin - Tailwind CSS Admin Dashboard Template",
                    component: <FormElements />
                },
                {
                    path: "/forms/form-layout",
                    title: "Form Layout | TailAdmin - Tailwind CSS Admin Dashboard Template",
                    component: <FormLayout />
                },
                {
                    path: "/tables",
                    title: "Tables | TailAdmin - Tailwind CSS Admin Dashboard Template",
                    component: <Tables />
                },
                {
                    path: "/settings",
                    title: "Settings | TailAdmin - Tailwind CSS Admin Dashboard Template",
                    component: <Settings />
                },
                {
                    path: "/chart",
                    title: "Basic Chart | TailAdmin - Tailwind CSS Admin Dashboard Template",
                    component: <Chart />
                },
                {
                    path: "/servers",
                    title: "Servers | R2DTools Control Panel",
                    name: "Servers",
                    component: <ServerList />
                },
                {
                    path: "/certificates",
                    title: "Certificates | R2DTools Control Panel",
                    name: "Certificates",
                    component: <CertificateList />
                },
                {
                    path: "/servers/:guid",
                    title: "Server | R2DTools Control Panel",
                    name: 'Server',
                    component: <Server />,
                },
                {
                    path: "/servers/:guid/domains",
                    title: "Domains | R2DTools Control Panel",
                    name: 'Domains',
                    component: <DomainList />,
                },
                {
                    path: "/servers/:guid/certificates",
                    title: "Certificates | R2DTools Control Panel",
                    name: 'Certificates',
                    component: <CertificateList />,
                },
                {
                    path: "/servers/:guid/domains/:name",
                    title: "Domain | R2DTools Control Panel",
                    name: 'Domain',
                    component: <Domain />,
                },
                {
                    path: "/ui/alerts",
                    title: "Alerts | TailAdmin - Tailwind CSS Admin Dashboard Template",
                    component: <Alerts />
                },
                {
                    path: "/ui/buttons",
                    title: "Buttons | TailAdmin - Tailwind CSS Admin Dashboard Template",
                    component: <Buttons />
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
                    title: "Signin | R2DTools - Tailwind CSS Admin Dashboard Template",
                    component: <SignIn onSignIn={handleSignIn} />
                },
                {
                    path: "/auth/signup",
                    public: true,
                    title: "Signup | R2DTools - Tailwind CSS Admin Dashboard Template",
                    component: <SignUp />
                },
            ],
        },
    ];

    let allRoutes: RouteItem[] = [];
    routes.forEach(({ routes }) => allRoutes = allRoutes.concat(routes));

    const authRoutes = ["/auth/signin", "/auth/signup"];

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

        if ((currentUser !== null && isAuthRoute) || currentRoute === null) {
            navigate("/");
        }
    }, [currentUser, routes, authRoutes]);

    const filterRoutes = (item: RouteItem) => (currentUser !== null && !authRoutes.includes(item.path || "")) || (currentUser === null && item.public);

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
