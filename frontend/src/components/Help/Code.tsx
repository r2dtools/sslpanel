import React from 'react';
import SyntaxHighlighter from 'react-syntax-highlighter';
import { nightOwl as hStyle } from 'react-syntax-highlighter/dist/esm/styles/hljs';

interface CodeProps {
    children: string | string[];
};

const Code: React.FC<CodeProps> = ({ children }) => {
    return (
        <SyntaxHighlighter language="bash" style={hStyle} customStyle={{ borderRadius: '4px' }}>
            {children}
        </SyntaxHighlighter>
    )
};

export default Code;
