'use client';

import { Documentation } from "./openrpc-doc";
import { OpenrpcDocument } from "./openrpc-doc/tool-types";
import uploadSchema from '../upload.openrpc.json';

export default function Home() {
  return <Documentation schema={uploadSchema as OpenrpcDocument} />;
}
