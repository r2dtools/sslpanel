import { FC } from 'react';
import { NavLink } from 'react-router-dom';

interface DropdownMenuItemProps {
    name: string;
    url: string;
}

const DropdownMenuItem: FC<DropdownMenuItemProps> = ({ name, url }) => {
    return (
        <li>
            <NavLink
                to={url}
                className={({ isActive }) => 'group relative flex items-center gap-2.5 rounded-md px-4 font-medium text-bodydark2 duration-300 ease-in-out hover:text-white ' + (isActive && '!text-white')}
            >
                {name}
            </NavLink>
        </li>
    )
};

export default DropdownMenuItem;
