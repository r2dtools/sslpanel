import { generatePath, Link, matchPath, useLocation, useParams } from 'react-router-dom';
import React, { useContext, useMemo } from 'react';
import RoutesContext from '../app/context';
import { RouteItem } from '../app/types';

interface BreadcrumbProps {
    pageName: string;
    children?: React.ReactNode;
}

const createBreadcrumbs = (pathname: string, routes: RouteItem[], params: {} | undefined): { name: string, url: string }[] => {
    const breadcrumbs = [];
    const parts = pathname.split('/');

    while (parts.length > 0) {
        parts.pop();
        const match = routes.find(
            (routeItem: RouteItem) => routeItem.path
                && routeItem.path !== '*'
                && matchPath(routeItem.path, parts.join('/')) !== null
        );

        if (match?.name && match.path) {
            breadcrumbs.unshift({
                name: match.name,
                url: generatePath(match.path, params),
            });
        }
    }

    return breadcrumbs;
};

const Breadcrumb = ({ pageName, children }: BreadcrumbProps) => {
    const routes = useContext(RoutesContext);
    const { pathname } = useLocation();
    const params = useParams();
    const breadcrumbs = useMemo(() => createBreadcrumbs(pathname, routes, params), [pathname, routes, params]);

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
                    {breadcrumbs.map(({ name, url }) => <Link className='font-medium' to={url} key={name}>{`${name} /`}</Link>)}
                    <li className="font-medium text-primary">{pageName}</li>
                </ol>
            </nav>
        </div>
    );
};

export default Breadcrumb;
