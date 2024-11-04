import { createContext } from "react";
import { User } from "./types";

const AuthContext = createContext<User | null>(null);

export default AuthContext;
