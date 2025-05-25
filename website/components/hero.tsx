import Link from "next/link";
import { motion } from "framer-motion";

export const meta = {
  title: "Goofer ORM",
  description:
    "A powerful, type safe ORM for Go that gives your structs a life of their own with relationships, migrations, and zero drama.",
};

export default function Hero() {
  return (
    <div className="h-full w-full flex items-center justify-center px-4 py-16 bg-gradient-to-br from-gray-900 via-black to-gray-800 text-white">
      <div className="text-center max-w-3xl">
        <motion.h1
          className="text-5xl font-extrabold mb-4 tracking-tight"
          initial={{ opacity: 0, y: -30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6 }}
        >
          {meta.title}
        </motion.h1>

        <motion.p
          className="italic text-lg mb-2 text-gray-300"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.3, duration: 0.5 }}
        >
          Stop hand writing SQL like it's 1999
        </motion.p>

        <motion.h2
          className="text-xl font-medium mb-8 text-gray-400"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.6, duration: 0.5 }}
        >
          {meta.description}
        </motion.h2>

        <motion.div
          className="flex justify-center gap-4"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.9, duration: 0.5 }}
        >
          <Link
            href="/docs"
            className="px-6 py-3 rounded-lg bg-blue-600 hover:bg-blue-700 transition font-semibold shadow-lg"
          >
            Read the docs <span className="ml-1">→</span>
          </Link>

          <Link
            href="https://github.com/gooferOrm/goofer"
            className="px-6 py-3 rounded-lg border border-gray-500 hover:border-white transition font-semibold text-white"
          >
            GitHub Repo <span className="ml-1">→</span>
          </Link>
        </motion.div>
      </div>
    </div>
  );
}