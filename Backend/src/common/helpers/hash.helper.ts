import * as bcrypt from "bcrypt"

export class hashing {
    private static readonly SALT = 10

    static async hash(password : string) {
        return bcrypt.hash(password, this.SALT)
    }

    static async compare(password : string, hashPassword : string) {
        return bcrypt.compare(password, hashPassword)
    }
}