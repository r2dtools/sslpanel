import centosIcon from '../../images/server/centos-icon.svg';
import debianIcon from '../../images/server/debian-icon.svg';
import ubuntuIcon from '../../images/server/ubuntu-icon.svg';
import linuxIcon from '../../images/server/linux-icon.svg';
import { OsCode } from './types';

export const getOsIcon = (code: string): string => {
    let icon = '';

    switch (code) {
        case OsCode.Ubuntu:
            icon = ubuntuIcon;

            break;
        case OsCode.Censos:
            icon = centosIcon;

            break;
        case OsCode.Debian:
            icon = debianIcon;

            break;
        default:
            icon = linuxIcon;
    }

    return icon;
};

export const getOsName = (code: string): string => {
    let name = '';

    switch (code) {
        case OsCode.Ubuntu:
            name = 'Ubuntu';

            break;
        case OsCode.Censos:
            name = 'CentOs';

            break;
        case OsCode.Debian:
            name = 'Debian';

            break;
        default:
            name = 'Linux';
    }

    return name;
};