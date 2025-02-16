import { Drawer } from "flowbite-react";
import { FC } from "react";
import { HiMiniCodeBracketSquare } from "react-icons/hi2";
import SyntaxHighlighter from 'react-syntax-highlighter';
import { dracula, a11yLight } from 'react-syntax-highlighter/dist/esm/styles/hljs';
import { Domain } from "../types";
import { ColorTheme } from "../../../types/theme";
import { useAppSelector } from "../../../app/hooks";
import { selectColorMode } from "../../../app/appSlice";

type DomainConfigDrawerProps = {
    domain: Domain;
    authToken: string;
    guid: string;
    open: boolean;
    config: string | null;
    onClose: () => void;
};

const DomainConfigDrawer: FC<DomainConfigDrawerProps> = ({ domain, open, onClose, config }) => {
    const appColorMode = useAppSelector(selectColorMode);

    const handleFormClose = (): void => {
        onClose();
    };

    return (
        <Drawer className='z-[999] min-w-[50%]' open={open} onClose={handleFormClose} position='right'>
            <Drawer.Header title={`Configuration for domain ${domain.servername}`} titleIcon={HiMiniCodeBracketSquare} />
            <Drawer.Items>
                {
                    config
                        ? (
                            <SyntaxHighlighter
                                language={domain.webserver}
                                style={appColorMode === ColorTheme.Dark ? dracula : a11yLight}
                                customStyle={{ "background": "none" }}
                            >
                                {config}
                            </SyntaxHighlighter>
                        )
                        : null
                }

            </Drawer.Items>
        </Drawer>
    );
};

export default DomainConfigDrawer;
