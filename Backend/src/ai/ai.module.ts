import { Module } from "@nestjs/common";
import { ConfigModule } from "@nestjs/config";
import { AiServiceClient } from "./ai-service.client";

@Module({
    imports : [ConfigModule],
    providers : [AiServiceClient],
    exports : [AiServiceClient]
})
export class AiModule {}