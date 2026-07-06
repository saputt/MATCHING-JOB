import { Controller, Param, Post, UseGuards } from "@nestjs/common";
import { AnalysisService } from "./analysis.service";
import { GetUser } from "src/common/decorators/get-user.decorator";
import { JwtAuthGuard } from "src/common/guards/jwt-auth.guard";

@Controller("analyze")
@UseGuards(JwtAuthGuard)
export class AnalysisController {
    constructor(private readonly service : AnalysisService) {}

    @Post(":jobId")
    async getOrCreateAnalysis(@GetUser('id') userId : string, @Param("jobId") jobId : string) {
        const res = await this.service.getOrCreateAnalysis(userId, jobId)
        return {
            message : "analysis success",
            data : res
        }
    }
}