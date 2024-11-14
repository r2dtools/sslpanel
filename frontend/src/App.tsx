import React, { useEffect } from 'react';
import { Route, Routes, useLocation, useNavigate } from 'react-router-dom';
import Loader from './common/Loader';
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
import { FetchStatus } from './app/types';
import { ColorTheme } from './types/theme';
import AuthLayout from './layout/AuthLayout';
import ServerList from './pages/ServerList';

interface RouteItem {
    path?: string;
    index?: boolean;
    public?: boolean;
    title: string;
    component: React.ReactNode | null;
};

function App() {
    const { pathname } = useLocation();
    const navigate = useNavigate();
    const [colorMode] = useColorMode();

    const currentUser = useAppSelector(selectCurrentUser);
    const currentUserLoadStatus = useAppSelector(selectCurrentUserFetchStatus);
    const [authToken, setAuthToken] = useLocalStorage<string | null>("r2panel-token", null);
    const [authTokenExpire, setAuthTokenExpire] = useLocalStorage<string | null>("r2panel-token-expire", null);
    const dispatch = useAppDispatch();

    const loading = currentUserLoadStatus === FetchStatus.Pending;

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
                    component: <ServerList />
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

        if (currentUser === null && !currentRoute?.public) {
            navigate("/auth/signin");

            return;
        }

        const isAuthRoute = authRoutes.includes(currentRoute?.path || "");

        if ((currentUser !== null && isAuthRoute) || currentRoute === null) {
            navigate("/");
        }
    }, [currentUser, routes, authRoutes]);

    return loading ? (
        <Loader />
    ) : (
        <AuthContext.Provider value={currentUser}>
            <Routes>
                {
                    routes.map(({ layout: Layout, routes, props }) => (
                        <Route element={<Layout {...props} />} key={Layout.name}>
                            {
                                routes
                                    .filter((item: RouteItem) => (currentUser !== null && !authRoutes.includes(item.path || "")) || (currentUser === null && item.public))
                                    .map((item: RouteItem) => (
                                        <Route
                                            key={item.path}
                                            index={item.index}
                                            path={item.path}
                                            element={
                                                <>
                                                    <PageTitle title={item.title} />
                                                    {item.component}
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
        </AuthContext.Provider>
    );
}

export default App;
