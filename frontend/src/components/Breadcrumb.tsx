import { Link } from 'react-router-dom';
import React from 'react';

interface BreadcrumbProps {
    pageName: string;
    children?: React.ReactNode;
}

const Breadcrumb = ({ pageName, children }: BreadcrumbProps) => {
    return (
        <div className="mb-6 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <div className='flex flex-row gap-5 items-center justify-between sm:justify-start'>
                <h2 className="text-title-md2 font-semibold text-black dark:text-white">
                    {pageName}
                </h2>
                {children}
            </div>

            <nav>
                <ol className="flex items-center gap-2">
                    <li>
                        <Link className="font-medium" to="/">
                            Dashboard /
                        </Link>
                    </li>
                    <li className="font-medium text-primary">{pageName}</li>
                </ol>
            </nav>
        </div>
    );
};

export default Breadcrumb;
