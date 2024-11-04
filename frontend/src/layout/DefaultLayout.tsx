import React, { useState, useContext } from 'react';
import { Outlet } from 'react-router-dom';
import Header from '../components/Header/index';
import Sidebar from '../components/Sidebar/index';
import AuthContext from '../features/auth/context';

interface DefaultLayoutProps {
    onSignOut: () => void;
};

const DefaultLayout: React.FC<DefaultLayoutProps> = ({ onSignOut }) => {
    const currentUser = useContext(AuthContext);
    const [sidebarOpen, setSidebarOpen] = useState(true);

    return (
        <div className="dark:bg-boxdark-2 dark:text-bodydark">
            {/* <!-- ===== Page Wrapper Start ===== --> */}
            <div className="flex h-screen overflow-hidden">
                {currentUser !== null && <Sidebar sidebarOpen={sidebarOpen} setSidebarOpen={setSidebarOpen} />}

                {/* <!-- ===== Content Area Start ===== --> */}
                <div className="relative flex flex-1 flex-col overflow-y-auto overflow-x-hidden">
                    {currentUser !== null && <Header sidebarOpen={sidebarOpen} setSidebarOpen={setSidebarOpen} onSignOut={onSignOut} />}

                    {/* <!-- ===== Main Content Start ===== --> */}
                    <main>
                        <div className="mx-auto max-w-screen-2xl p-4 md:p-6 2xl:p-10">
                            <Outlet />
                        </div>
                    </main>
                    {/* <!-- ===== Main Content End ===== --> */}
                </div>
                {/* <!-- ===== Content Area End ===== --> */}
            </div>
            {/* <!-- ===== Page Wrapper End ===== --> */}
        </div>
    );
};

export default DefaultLayout;
